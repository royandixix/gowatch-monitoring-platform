package monitor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/royandi/gowatch/backend/internal/model"
)

const (
	maxTargetURLLength = 2048
	maxRedirects       = 5
)

var blockedNetworks = mustParseNetworks(
	[]string{
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/4",
		"240.0.0.0/4",

		"::/128",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
		"ff00::/8",
		"2001:db8::/32",
	},
)

var secureHTTPClient = newSecureHTTPClient()

func mustParseNetworks(
	values []string,
) []*net.IPNet {
	networks := make(
		[]*net.IPNet,
		0,
		len(values),
	)

	for _, value := range values {
		_, network, err :=
			net.ParseCIDR(value)

		if err != nil {
			panic(
				"CIDR keamanan tidak valid: " +
					value,
			)
		}

		networks = append(
			networks,
			network,
		)
	}

	return networks
}

func isBlockedIP(
	ip net.IP,
) bool {
	if ip == nil {
		return true
	}

	if ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() {
		return true
	}

	for _, network := range blockedNetworks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

func validateStaticHostname(
	hostname string,
) error {
	hostname = strings.ToLower(
		strings.TrimSuffix(
			strings.TrimSpace(hostname),
			".",
		),
	)

	if hostname == "" {
		return errors.New(
			"hostname wajib diisi",
		)
	}

	blockedNames :=
		[]string{
			"localhost",
			"localhost.localdomain",
		}

	for _, blockedName := range blockedNames {
		if hostname == blockedName {
			return errors.New(
				"alamat lokal tidak diizinkan",
			)
		}
	}

	blockedSuffixes :=
		[]string{
			".localhost",
			".local",
			".internal",
			".lan",
			".home",
		}

	for _, suffix := range blockedSuffixes {
		if strings.HasSuffix(
			hostname,
			suffix,
		) {
			return errors.New(
				"hostname jaringan internal tidak diizinkan",
			)
		}
	}

	if ip :=
		net.ParseIP(hostname); ip != nil {
		if isBlockedIP(ip) {
			return errors.New(
				"alamat IP privat atau internal tidak diizinkan",
			)
		}

		return nil
	}

	// Host tunggal seperti "redis", "database",
	// "backend", dan sejenisnya tidak diizinkan.
	if !strings.Contains(
		hostname,
		".",
	) {
		return errors.New(
			"hostname lokal tidak diizinkan",
		)
	}

	return nil
}

func parseAndNormalizeURL(
	rawURL string,
) (string, *url.URL, error) {
	rawURL =
		strings.TrimSpace(
			rawURL,
		)

	if rawURL == "" {
		return "",
			nil,
			errors.New(
				"URL wajib diisi",
			)
	}

	if len(rawURL) >
		maxTargetURLLength {
		return "",
			nil,
			errors.New(
				"URL terlalu panjang",
			)
	}

	parsedURL, err :=
		url.Parse(rawURL)

	if err != nil {
		return "",
			nil,
			errors.New(
				"format URL tidak valid",
			)
	}

	parsedURL.Scheme =
		strings.ToLower(
			parsedURL.Scheme,
		)

	if parsedURL.Scheme != "http" &&
		parsedURL.Scheme != "https" {
		return "",
			nil,
			errors.New(
				"URL hanya boleh menggunakan http atau https",
			)
	}

	if parsedURL.Host == "" {
		return "",
			nil,
			errors.New(
				"hostname URL wajib diisi",
			)
	}

	if parsedURL.User != nil {
		return "",
			nil,
			errors.New(
				"credential pada URL tidak diizinkan",
			)
	}

	hostname :=
		strings.ToLower(
			strings.TrimSuffix(
				parsedURL.Hostname(),
				".",
			),
		)

	if err :=
		validateStaticHostname(
			hostname,
		); err != nil {
		return "",
			nil,
			err
	}

	port :=
		parsedURL.Port()

	if port != "" {
		portNumber, err :=
			strconv.Atoi(port)

		if err != nil ||
			portNumber < 1 ||
			portNumber > 65535 {
			return "",
				nil,
				errors.New(
					"port URL tidak valid",
				)
		}
	}

	if ip :=
		net.ParseIP(hostname); ip != nil {
		if port != "" {
			parsedURL.Host =
				net.JoinHostPort(
					hostname,
					port,
				)
		} else if strings.Contains(
			hostname,
			":",
		) {
			parsedURL.Host =
				"[" +
					hostname +
					"]"
		} else {
			parsedURL.Host =
				hostname
		}
	} else {
		if port != "" {
			parsedURL.Host =
				net.JoinHostPort(
					hostname,
					port,
				)
		} else {
			parsedURL.Host =
				hostname
		}
	}

	// Fragment tidak pernah dikirim ke server tujuan.
	parsedURL.Fragment = ""
	parsedURL.RawFragment = ""

	return parsedURL.String(),
		parsedURL,
		nil
}

func resolvePublicIPs(
	ctx context.Context,
	hostname string,
) ([]net.IP, error) {
	if err :=
		validateStaticHostname(
			hostname,
		); err != nil {
		return nil, err
	}

	if ip :=
		net.ParseIP(hostname); ip != nil {
		if isBlockedIP(ip) {
			return nil,
				errors.New(
					"alamat IP privat atau internal tidak diizinkan",
				)
		}

		return []net.IP{
			ip,
		}, nil
	}

	resolveContext, cancel :=
		context.WithTimeout(
			ctx,
			3*time.Second,
		)

	defer cancel()

	addresses, err :=
		net.DefaultResolver.LookupIPAddr(
			resolveContext,
			hostname,
		)

	if err != nil {
		return nil,
			errors.New(
				"hostname tidak dapat ditemukan",
			)
	}

	if len(addresses) == 0 {
		return nil,
			errors.New(
				"hostname tidak memiliki alamat IP",
			)
	}

	publicIPs :=
		make(
			[]net.IP,
			0,
			len(addresses),
		)

	for _, address := range addresses {
		if isBlockedIP(
			address.IP,
		) {
			return nil,
				errors.New(
					"hostname mengarah ke jaringan privat atau internal",
				)
		}

		publicIPs = append(
			publicIPs,
			address.IP,
		)
	}

	if len(publicIPs) == 0 {
		return nil,
			errors.New(
				"target tidak memiliki alamat IP publik",
			)
	}

	return publicIPs, nil
}

// NormalizeAndValidateURL memvalidasi URL,
// melakukan DNS resolution, dan memastikan target
// tidak mengarah ke jaringan lokal/private.
func NormalizeAndValidateURL(
	ctx context.Context,
	rawURL string,
) (string, error) {
	normalizedURL,
		parsedURL,
		err :=
		parseAndNormalizeURL(
			rawURL,
		)

	if err != nil {
		return "", err
	}

	_, err =
		resolvePublicIPs(
			ctx,
			parsedURL.Hostname(),
		)

	if err != nil {
		return "", err
	}

	return normalizedURL, nil
}

// IsValidURL dipertahankan untuk kompatibilitas.
// Fungsi ini hanya melakukan validasi sintaks
// dan pemeriksaan target statis.
func IsValidURL(
	targetURL string,
) bool {
	_, _, err :=
		parseAndNormalizeURL(
			targetURL,
		)

	return err == nil
}

func secureDialContext(
	ctx context.Context,
	network string,
	address string,
) (net.Conn, error) {
	hostname,
		port,
		err :=
		net.SplitHostPort(
			address,
		)

	if err != nil {
		return nil,
			fmt.Errorf(
				"alamat target tidak valid: %w",
				err,
			)
	}

	hostname =
		strings.Trim(
			hostname,
			"[]",
		)

	publicIPs, err :=
		resolvePublicIPs(
			ctx,
			hostname,
		)

	if err != nil {
		return nil,
			fmt.Errorf(
				"target diblokir: %w",
				err,
			)
	}

	dialer :=
		&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}

	var lastError error

	for _, ip := range publicIPs {
		target :=
			net.JoinHostPort(
				ip.String(),
				port,
			)

		connection, err :=
			dialer.DialContext(
				ctx,
				network,
				target,
			)

		if err == nil {
			return connection, nil
		}

		lastError = err
	}

	if lastError == nil {
		lastError =
			errors.New(
				"tidak ada alamat publik yang dapat dihubungi",
			)
	}

	return nil, lastError
}

func newSecureHTTPClient() *http.Client {
	transport :=
		&http.Transport{
			// Proxy environment tidak digunakan.
			// Target selalu diperiksa langsung.
			Proxy: nil,

			DialContext: secureDialContext,

			ForceAttemptHTTP2: true,

			MaxIdleConns: 100,

			MaxIdleConnsPerHost: 10,

			IdleConnTimeout: 30 * time.Second,

			TLSHandshakeTimeout: 5 * time.Second,

			ResponseHeaderTimeout: 8 * time.Second,

			ExpectContinueTimeout: 1 * time.Second,
		}

	return &http.Client{
		Transport: transport,

		Timeout: 10 * time.Second,

		CheckRedirect: func(
			request *http.Request,
			via []*http.Request,
		) error {
			if len(via) >=
				maxRedirects {
				return errors.New(
					"terlalu banyak redirect",
				)
			}

			_, parsedURL, err :=
				parseAndNormalizeURL(
					request.URL.String(),
				)

			if err != nil {
				return fmt.Errorf(
					"redirect tidak diizinkan: %w",
					err,
				)
			}

			_, err =
				resolvePublicIPs(
					request.Context(),
					parsedURL.Hostname(),
				)

			if err != nil {
				return fmt.Errorf(
					"redirect diblokir: %w",
					err,
				)
			}

			return nil
		},
	}
}

// CheckWebsiteContext mengecek website/API
// menggunakan context dan HTTP client yang aman.
func CheckWebsiteContext(
	ctx context.Context,
	targetURL string,
) model.MonitorResult {
	start :=
		time.Now()

	normalizedURL, err :=
		NormalizeAndValidateURL(
			ctx,
			targetURL,
		)

	if err != nil {
		return model.MonitorResult{
			URL: targetURL,

			Status: "DOWN",

			HTTPStatus: 0,

			ResponseTime: time.Since(
				start,
			).Milliseconds(),

			CheckedAt: time.Now().
				Format(
					time.RFC3339,
				),

			Error: err.Error(),
		}
	}

	request, err :=
		http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			normalizedURL,
			nil,
		)

	if err != nil {
		return model.MonitorResult{
			URL: normalizedURL,

			Status: "DOWN",

			ResponseTime: time.Since(
				start,
			).Milliseconds(),

			CheckedAt: time.Now().
				Format(
					time.RFC3339,
				),

			Error: "gagal membuat HTTP request",
		}
	}

	request.Header.Set(
		"User-Agent",
		"GoWatch-Monitor/1.0",
	)

	request.Header.Set(
		"Accept",
		"*/*",
	)

	response, err :=
		secureHTTPClient.Do(
			request,
		)

	responseTime :=
		time.Since(
			start,
		).Milliseconds()

	result :=
		model.MonitorResult{
			URL: normalizedURL,

			Status: "DOWN",

			HTTPStatus: 0,

			ResponseTime: responseTime,

			CheckedAt: time.Now().
				Format(
					time.RFC3339,
				),
		}

	if err != nil {
		result.Error =
			err.Error()

		return result
	}

	defer response.Body.Close()

	result.HTTPStatus =
		response.StatusCode

	if response.StatusCode >=
		200 &&
		response.StatusCode <
			400 {
		result.Status =
			"UP"
	}

	return result
}

// CheckWebsite dipertahankan agar kode lama tetap kompatibel.
func CheckWebsite(
	targetURL string,
) model.MonitorResult {
	return CheckWebsiteContext(
		context.Background(),
		targetURL,
	)
}

func CheckWebsitesConcurrentlyContext(
	ctx context.Context,
	targetURLs []string,
) []model.MonitorResult {
	results :=
		make(
			[]model.MonitorResult,
			len(targetURLs),
		)

	var waitGroup sync.WaitGroup

	for index, targetURL := range targetURLs {
		waitGroup.Add(1)

		go func(
			resultIndex int,
			currentURL string,
		) {
			defer waitGroup.Done()

			results[resultIndex] =
				CheckWebsiteContext(
					ctx,
					currentURL,
				)
		}(
			index,
			targetURL,
		)
	}

	waitGroup.Wait()

	return results
}

func CheckWebsitesConcurrently(
	targetURLs []string,
) []model.MonitorResult {
	return CheckWebsitesConcurrentlyContext(
		context.Background(),
		targetURLs,
	)
}
