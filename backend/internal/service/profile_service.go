package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"github.com/royandi/gowatch/backend/internal/model"
	"github.com/royandi/gowatch/backend/internal/repository"
)

var (
	ErrInvalidProfileName = errors.New(
		"nama minimal 2 karakter",
	)

	ErrCurrentPasswordInvalid = errors.New(
		"password lama tidak sesuai",
	)

	ErrNewPasswordSame = errors.New(
		"password baru harus berbeda dari password lama",
	)
)

type ProfileService struct {
	userRepository *repository.UserRepository
}

func NewProfileService(
	userRepository *repository.UserRepository,
) *ProfileService {
	return &ProfileService{
		userRepository: userRepository,
	}
}

func (s *ProfileService) UpdateProfile(
	ctx context.Context,
	userID int64,
	request model.UpdateProfileRequest,
) (model.User, error) {
	name :=
		strings.TrimSpace(
			request.Name,
		)

	email :=
		strings.ToLower(
			strings.TrimSpace(
				request.Email,
			),
		)

	if len([]rune(name)) < 2 {
		return model.User{},
			ErrInvalidProfileName
	}

	user, err :=
		s.userRepository.UpdateProfile(
			ctx,
			userID,
			name,
			email,
		)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(
			err,
			&pgErr,
		) &&
			pgErr.Code == "23505" {
			return model.User{},
				ErrEmailExists
		}

		return model.User{}, err
	}

	return user, nil
}

func (s *ProfileService) ChangePassword(
	ctx context.Context,
	userID int64,
	request model.ChangePasswordRequest,
) error {
	user, err :=
		s.userRepository.FindByID(
			ctx,
			userID,
		)

	if err != nil {
		return err
	}

	err =
		bcrypt.CompareHashAndPassword(
			[]byte(
				user.PasswordHash,
			),
			[]byte(
				request.CurrentPassword,
			),
		)

	if err != nil {
		return ErrCurrentPasswordInvalid
	}

	if request.CurrentPassword ==
		request.NewPassword {
		return ErrNewPasswordSame
	}

	passwordHash, err :=
		bcrypt.GenerateFromPassword(
			[]byte(
				request.NewPassword,
			),
			bcrypt.DefaultCost,
		)

	if err != nil {
		return err
	}

	return s.userRepository.UpdatePassword(
		ctx,
		userID,
		string(passwordHash),
	)
}
