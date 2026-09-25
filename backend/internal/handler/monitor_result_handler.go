package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/royandi/gowatch/backend/internal/model"
	"github.com/royandi/gowatch/backend/internal/repository"
)

type MonitorResultHandler struct {
	monitorResultRepository *repository.MonitorResultRepository
}

func NewMonitorResultHandler(
	monitorResultRepository *repository.MonitorResultRepository,
) *MonitorResultHandler {
	return &MonitorResultHandler{
		monitorResultRepository: monitorResultRepository,
	}
}

func getPositiveIntQuery(
	c *gin.Context,
	key string,
	defaultValue int,
	maxValue int,
) (int, bool) {
	value :=
		strings.TrimSpace(
			c.Query(key),
		)

	if value == "" {
		return defaultValue, true
	}

	parsed, err :=
		strconv.Atoi(
			value,
		)

	if err != nil ||
		parsed < 1 {
		return 0, false
	}

	if maxValue > 0 &&
		parsed > maxValue {
		parsed = maxValue
	}

	return parsed, true
}

func getResultLimit(
	c *gin.Context,
) int {
	limit, ok :=
		getPositiveIntQuery(
			c,
			"limit",
			100,
			200,
		)

	if !ok {
		return 100
	}

	return limit
}

func (h *MonitorResultHandler) History(
	c *gin.Context,
) {
	userID, ok :=
		getUserID(c)

	if !ok {
		return
	}

	page, valid :=
		getPositiveIntQuery(
			c,
			"page",
			1,
			0,
		)

	if !valid {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "page tidak valid",
			},
		)

		return
	}

	perPageValue :=
		strings.TrimSpace(
			c.Query(
				"per_page",
			),
		)

	if perPageValue == "" {
		perPageValue =
			strings.TrimSpace(
				c.Query(
					"limit",
				),
			)
	}

	perPage := 20

	if perPageValue != "" {
		parsed, err :=
			strconv.Atoi(
				perPageValue,
			)

		if err != nil ||
			parsed < 1 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "per_page tidak valid",
				},
			)

			return
		}

		perPage = parsed
	}

	if perPage > 100 {
		perPage = 100
	}

	status :=
		strings.ToUpper(
			strings.TrimSpace(
				c.Query(
					"status",
				),
			),
		)

	if status == "ALL" {
		status = ""
	}

	if status != "" &&
		status != "UP" &&
		status != "DOWN" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "status hanya boleh ALL, UP, atau DOWN",
			},
		)

		return
	}

	var monitorID int64

	monitorIDValue :=
		strings.TrimSpace(
			c.Query(
				"monitor_id",
			),
		)

	if monitorIDValue != "" {
		parsed, err :=
			strconv.ParseInt(
				monitorIDValue,
				10,
				64,
			)

		if err != nil ||
			parsed < 1 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "monitor_id tidak valid",
				},
			)

			return
		}

		monitorID = parsed
	}

	results, total, err :=
		h.monitorResultRepository.FindHistoryByUser(
			c.Request.Context(),
			userID,
			model.MonitorHistoryFilter{
				Status:    status,
				MonitorID: monitorID,
				Page:      page,
				PerPage:   perPage,
			},
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "gagal mengambil history monitoring",
			},
		)

		return
	}

	totalPages := 0

	if total > 0 {
		totalPages =
			int(
				(total +
					int64(perPage) -
					1) /
					int64(perPage),
			)
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"total":       total,
			"page":        page,
			"per_page":    perPage,
			"total_pages": totalPages,
			"results":     results,
		},
	)
}

func (h *MonitorResultHandler) ByMonitor(
	c *gin.Context,
) {
	userID, ok :=
		getUserID(c)

	if !ok {
		return
	}

	monitorID, ok :=
		getMonitorID(c)

	if !ok {
		return
	}

	limit :=
		getResultLimit(c)

	results, err :=
		h.monitorResultRepository.FindByMonitorAndUser(
			c.Request.Context(),
			monitorID,
			userID,
			limit,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "gagal mengambil hasil monitoring",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"total":   len(results),
			"results": results,
		},
	)
}

func (h *MonitorResultHandler) Stats(
	c *gin.Context,
) {
	userID, ok :=
		getUserID(c)

	if !ok {
		return
	}

	monitorID, ok :=
		getMonitorID(c)

	if !ok {
		return
	}

	stats, err :=
		h.monitorResultRepository.GetStatsByMonitorAndUser(
			c.Request.Context(),
			monitorID,
			userID,
		)

	if errors.Is(
		err,
		pgx.ErrNoRows,
	) {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "monitor tidak ditemukan",
			},
		)

		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "gagal mengambil statistik monitor",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"stats": stats,
		},
	)
}
