package user

import (
	"context"
	"errors"
	"time"

	"taskflow/internal/config"

	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetAllUsers(ctx context.Context, search string) ([]User, error)
}

type service struct {
	repo Repository
	cfg  *config.Config
	log  *slog.Logger
}

func NewService(repo Repository, cfg *config.Config, log *slog.Logger) Service {
	return &service{repo: repo, cfg: cfg, log: log}
}

func (s *service) Register(ctx context.Context, name, email, password string) (string, error) {
	s.log.Info("Registering new user", "email", email, "name", name)

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		s.log.Error("Failed to hash password", "error", err)
		return "", err
	}

	user := &User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		s.log.Error("Failed to create user", "error", err, "email", email)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		s.log.Error("Failed to sign JWT", "error", err, "email", email)
		return "", err
	}

	s.log.Info("User logged in successfully", "email", email)
	return tokenString, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	s.log.Info("User login attempt", "email", email)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Warn("Login failed: user not found", "email", email)
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.log.Warn("Login failed: invalid password", "email", email)
		return "", errors.New("invalid credentials")
	}

	// create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		s.log.Error("Failed to sign JWT", "error", err, "email", email)
		return "", err
	}

	s.log.Info("User logged in successfully", "email", email)
	return tokenString, nil
}

func (s *service) GetAllUsers(ctx context.Context, search string) ([]User, error) {
	s.log.Info("Fetching all users", "search", search)

	users, err := s.repo.GetAllUsers(ctx, search)
	if err != nil {
		s.log.Error("Failed to fetch users", "error", err, "search", search)
		return nil, err
	}

	s.log.Info("Users fetched successfully", "count", len(users))
	return users, nil
}
