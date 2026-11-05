package services

import (
	"context"
	"course-management/models"
	"course-management/repository"
	"course-management/utils"
	"errors"
)

// AuthService defines authentication business logic
type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) error
	Login(ctx context.Context, req *models.LoginRequest) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) error {
	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Default role to LEARNER
	role := req.Role
	if role == "" {
		role = "LEARNER"
	}
	if role != "ADMIN" && role != "LEARNER" {
		return errors.New("invalid role, must be ADMIN or LEARNER")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	_, err = s.userRepo.Create(ctx, user)
	return err
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.UserID, user.Role)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
