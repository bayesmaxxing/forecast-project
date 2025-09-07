package services

import (
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  repository.UserRepository
	cache *cache.Cache
}

func NewUserService(repo repository.UserRepository, cache *cache.Cache) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Set the hashed password
	user.Password = string(hashedPassword)

	s.cache.Delete("users")
	// Create the user with the hashed password
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	s.cache.Delete("users")
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) ValidateUser(ctx context.Context, id int64) (bool, error) {
	return s.repo.ValidateUser(ctx, id)
}

func (s *UserService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// Get the user to verify old password
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// AdminResetPassword resets a user's password without requiring the old password
func (s *UserService) AdminResetPassword(ctx context.Context, userID int64, newPassword string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

func (s *UserService) VerifyPassword(ctx context.Context, username string, password string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	users_cache, found := s.cache.Get("users")
	if found {
		return users_cache.([]*models.User), nil
	}
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.Set("users", users)
	return users, nil
}
