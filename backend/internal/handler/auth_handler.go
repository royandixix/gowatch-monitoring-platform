package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/royandi/gowatch/backend/internal/model"
	"github.com/royandi/gowatch/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(
	authService *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(
	c *gin.Context,
) {
	var request model.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "data registrasi tidak valid",
		})

		return
	}

	user, err := h.authService.Register(
		c.Request.Context(),
		request,
	)

	if errors.Is(
		err,
		service.ErrEmailExists,
	) {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal membuat akun",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registrasi berhasil",
		"user":    user,
	})
}

func (h *AuthHandler) Login(
	c *gin.Context,
) {
	var request model.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "data login tidak valid",
		})

		return
	}

	user, token, err := h.authService.Login(
		c.Request.Context(),
		request,
	)

	if errors.Is(
		err,
		service.ErrInvalidCredentials,
	) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "login gagal",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
		"token":   token,
		"user":    user,
	})
}

func (h *AuthHandler) Me(
	c *gin.Context,
) {
	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user tidak ditemukan pada token",
		})

		return
	}

	userID, ok := userIDValue.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id tidak valid",
		})

		return
	}

	user, err := h.authService.GetUserByID(
		c.Request.Context(),
		userID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user tidak ditemukan",
		})

		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal mengambil data user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
