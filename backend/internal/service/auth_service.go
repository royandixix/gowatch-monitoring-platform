package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"github.com/royandi/gowatch/backend/internal/auth"
	"github.com/royandi/gowatch/backend/internal/model"
	"github.com/royandi/gowatch/backend/internal/repository"
)

var (
	ErrEmailExists        = errors.New("email sudah terdaftar")
	ErrInvalidCredentials = errors.New("email atau password salah")
)

type AuthService struct {
	userRepository *repository.UserRepository
	jwtSecret      string
}

func NewAuthService(
	userRepository *repository.UserRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	request model.RegisterRequest,
) (model.User, error) {
	name := strings.TrimSpace(request.Name)
	email := strings.ToLower(strings.TrimSpace(request.Email))

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return model.User{}, err
	}

	user, err := s.userRepository.Create(
		ctx,
		name,
		email,
		string(passwordHash),
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, ErrEmailExists
		}

		return model.User{}, err
	}

	return user, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	request model.LoginRequest,
) (model.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	user, err := s.userRepository.FindByEmail(
		ctx,
		email,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, "", ErrInvalidCredentials
	}

	if err != nil {
		return model.User{}, "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	)

	if err != nil {
		return model.User{}, "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(
		user.ID,
		user.Email,
		s.jwtSecret,
	)

	if err != nil {
		return model.User{}, "", err
	}

	return user, token, nil
}

func (s *AuthService) GetUserByID(
	ctx context.Context,
	id int64,
) (model.User, error) {
	return s.userRepository.FindByID(ctx, id)
}
