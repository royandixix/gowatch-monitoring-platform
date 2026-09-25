package monitor

import (
	"context"
	"strings"
	"testing"
)

func TestNormalizeAndValidateURL_PublicIP(
	t *testing.T,
) {
	t.Parallel()

	value, err :=
		NormalizeAndValidateURL(
			context.Background(),
			"https://1.1.1.1/test#section",
		)

	if err != nil {
		t.Fatalf(
			"URL publik seharusnya valid: %v",
			err,
		)
	}

	expected :=
		"https://1.1.1.1/test"

	if value != expected {
		t.Fatalf(
			"hasil normalisasi salah: got %q, want %q",
			value,
			expected,
		)
	}
}

func TestNormalizeAndValidateURL_BlockedTargets(
	t *testing.T,
) {
	t.Parallel()

	tests :=
		[]struct {
			name string
			url  string
		}{
			{
				name: "localhost",
				url:  "http://localhost:8081",
			},
			{
				name: "loopback IPv4",
				url:  "http://127.0.0.1:8081",
			},
			{
				name: "private 10",
				url:  "http://10.0.0.1",
			},
			{
				name: "private 172",
				url:  "http://172.16.0.1",
			},
			{
				name: "private 192",
				url:  "http://192.168.1.1",
			},
			{
				name: "metadata",
				url:  "http://169.254.169.254",
			},
			{
				name: "IPv6 loopback",
				url:  "http://[::1]",
			},
			{
				name: "local domain",
				url:  "http://server.local",
			},
			{
				name: "internal domain",
				url:  "http://database.internal",
			},
		}

	for _, test := range tests {
		test :=
			test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				_, err :=
					NormalizeAndValidateURL(
						context.Background(),
						test.url,
					)

				if err == nil {
					t.Fatalf(
						"URL seharusnya diblokir: %s",
						test.url,
					)
				}
			},
		)
	}
}

func TestNormalizeAndValidateURL_InvalidScheme(
	t *testing.T,
) {
	t.Parallel()

	tests :=
		[]string{
			"ftp://1.1.1.1",
			"file:///etc/passwd",
			"javascript:alert(1)",
		}

	for _, target := range tests {
		_, err :=
			NormalizeAndValidateURL(
				context.Background(),
				target,
			)

		if err == nil {
			t.Fatalf(
				"scheme harus ditolak: %s",
				target,
			)
		}
	}
}

func TestNormalizeAndValidateURL_RejectsCredentials(
	t *testing.T,
) {
	t.Parallel()

	_, err :=
		NormalizeAndValidateURL(
			context.Background(),
			"http://admin:password@1.1.1.1",
		)

	if err == nil {
		t.Fatal(
			"URL dengan credential seharusnya ditolak",
		)
	}
}

func TestIsValidURL(
	t *testing.T,
) {
	t.Parallel()

	tests :=
		[]struct {
			url      string
			expected bool
		}{
			{
				url:      "https://1.1.1.1",
				expected: true,
			},
			{
				url:      "http://127.0.0.1",
				expected: false,
			},
			{
				url:      "ftp://1.1.1.1",
				expected: false,
			},
			{
				url:      "",
				expected: false,
			},
		}

	for _, test := range tests {
		actual :=
			IsValidURL(
				test.url,
			)

		if actual !=
			test.expected {
			t.Fatalf(
				"IsValidURL(%q) = %v, want %v",
				test.url,
				actual,
				test.expected,
			)
		}
	}
}

func TestCheckWebsiteContext_BlocksPrivateTarget(
	t *testing.T,
) {
	t.Parallel()

	result :=
		CheckWebsiteContext(
			context.Background(),
			"http://127.0.0.1:8081",
		)

	if result.Status != "DOWN" {
		t.Fatalf(
			"status = %q, want DOWN",
			result.Status,
		)
	}

	if result.HTTPStatus != 0 {
		t.Fatalf(
			"http_status = %d, want 0",
			result.HTTPStatus,
		)
	}

	if result.Error == "" {
		t.Fatal(
			"error seharusnya terisi",
		)
	}

	if !strings.Contains(
		strings.ToLower(
			result.Error,
		),
		"internal",
	) &&
		!strings.Contains(
			strings.ToLower(
				result.Error,
			),
			"privat",
		) {
		t.Fatalf(
			"pesan error tidak sesuai: %q",
			result.Error,
		)
	}
}
