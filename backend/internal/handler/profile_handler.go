package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/royandi/gowatch/backend/internal/model"
	"github.com/royandi/gowatch/backend/internal/service"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(
	profileService *service.ProfileService,
) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (h *ProfileHandler) Update(
	c *gin.Context,
) {
	userID, ok :=
		getUserID(c)

	if !ok {
		return
	}

	var request model.UpdateProfileRequest

	if err :=
		c.ShouldBindJSON(
			&request,
		); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "data profile tidak valid",
			},
		)

		return
	}

	user, err :=
		h.profileService.UpdateProfile(
			c.Request.Context(),
			userID,
			request,
		)

	if err != nil {
		switch {
		case errors.Is(
			err,
			service.ErrEmailExists,
		):
			c.JSON(
				http.StatusConflict,
				gin.H{
					"error": err.Error(),
				},
			)

		case errors.Is(
			err,
			service.ErrInvalidProfileName,
		):
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)

		case errors.Is(
			err,
			pgx.ErrNoRows,
		):
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": "user tidak ditemukan",
				},
			)

		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "gagal memperbarui profile",
				},
			)
		}

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "profile berhasil diperbarui",
			"user":    user,
		},
	)
}

func (h *ProfileHandler) ChangePassword(
	c *gin.Context,
) {
	userID, ok :=
		getUserID(c)

	if !ok {
		return
	}

	var request model.ChangePasswordRequest

	if err :=
		c.ShouldBindJSON(
			&request,
		); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "data password tidak valid",
			},
		)

		return
	}

	err :=
		h.profileService.ChangePassword(
			c.Request.Context(),
			userID,
			request,
		)

	if err != nil {
		switch {
		case errors.Is(
			err,
			service.ErrCurrentPasswordInvalid,
		),
			errors.Is(
				err,
				service.ErrNewPasswordSame,
			):

			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)

		case errors.Is(
			err,
			pgx.ErrNoRows,
		):
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": "user tidak ditemukan",
				},
			)

		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "gagal mengubah password",
				},
			)
		}

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "password berhasil diubah",
		},
	)
}
