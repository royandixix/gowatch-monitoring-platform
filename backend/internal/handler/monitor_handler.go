package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/royandi/gowatch/backend/internal/model"
	"github.com/royandi/gowatch/backend/internal/service"
)

type MonitorHandler struct {
	monitorService *service.MonitorService
}

func NewMonitorHandler(
	monitorService *service.MonitorService,
) *MonitorHandler {
	return &MonitorHandler{
		monitorService: monitorService,
	}
}

func getUserID(
	c *gin.Context,
) (int64, bool) {
	value, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user tidak ditemukan pada token",
		})

		return 0, false
	}

	userID, ok := value.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id tidak valid",
		})

		return 0, false
	}

	return userID, true
}

func getMonitorID(
	c *gin.Context,
) (int64, bool) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id monitor tidak valid",
		})

		return 0, false
	}

	return id, true
}

func (h *MonitorHandler) Create(
	c *gin.Context,
) {
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	var request model.CreateMonitorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "data monitor tidak valid",
		})

		return
	}

	monitor, err := h.monitorService.Create(
		c.Request.Context(),
		userID,
		request,
	)

	if err != nil {
		switch {
		case errors.Is(
			err,
			service.ErrInvalidMonitorURL,
		):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(
			err,
			service.ErrInvalidMonitorInterval,
		):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(
			err,
			service.ErrInvalidMonitorType,
		):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal membuat monitor",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "monitor berhasil dibuat",
		"monitor": monitor,
	})
}

func (h *MonitorHandler) List(
	c *gin.Context,
) {
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	monitors, err := h.monitorService.List(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal mengambil monitor",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    len(monitors),
		"monitors": monitors,
	})
}

func (h *MonitorHandler) GetByID(
	c *gin.Context,
) {
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	id, ok := getMonitorID(c)

	if !ok {
		return
	}

	monitor, err := h.monitorService.GetByID(
		c.Request.Context(),
		id,
		userID,
	)

	if errors.Is(
		err,
		service.ErrMonitorNotFound,
	) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal mengambil monitor",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"monitor": monitor,
	})
}

func (h *MonitorHandler) Update(
	c *gin.Context,
) {
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	id, ok := getMonitorID(c)

	if !ok {
		return
	}

	var request model.UpdateMonitorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "data monitor tidak valid",
		})

		return
	}

	monitor, err := h.monitorService.Update(
		c.Request.Context(),
		id,
		userID,
		request,
	)

	if err != nil {
		switch {
		case errors.Is(
			err,
			service.ErrMonitorNotFound,
		):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(
			err,
			service.ErrInvalidMonitorURL,
		),
			errors.Is(
				err,
				service.ErrInvalidMonitorInterval,
			),
			errors.Is(
				err,
				service.ErrInvalidMonitorType,
			):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal memperbarui monitor",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "monitor berhasil diperbarui",
		"monitor": monitor,
	})
}

func (h *MonitorHandler) Delete(
	c *gin.Context,
) {
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	id, ok := getMonitorID(c)

	if !ok {
		return
	}

	err := h.monitorService.Delete(
		c.Request.Context(),
		id,
		userID,
	)

	if errors.Is(
		err,
		service.ErrMonitorNotFound,
	) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal menghapus monitor",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "monitor berhasil dihapus",
	})
}
