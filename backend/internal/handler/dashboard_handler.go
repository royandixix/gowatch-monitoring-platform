package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/royandi/gowatch/backend/internal/repository"
)

type DashboardHandler struct {
	dashboardRepository *repository.DashboardRepository
}

func NewDashboardHandler(
	dashboardRepository *repository.DashboardRepository,
) *DashboardHandler {
	return &DashboardHandler{
		dashboardRepository: dashboardRepository,
	}
}

func (h *DashboardHandler) Show(
	c *gin.Context,
) {
	userID, ok :=
		getUserID(c)

	if !ok {
		return
	}

	dashboard, err :=
		h.dashboardRepository.GetDashboard(
			c.Request.Context(),
			userID,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "gagal mengambil data dashboard",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		dashboard,
	)
}
