package service

import (
	"context"
	"fmt"

	"github.com/jwald3/go_rest_template/internal/domain"
	"github.com/jwald3/go_rest_template/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	Get(ctx context.Context, id int64) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(ctx context.Context, user *domain.User) error {
	// exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	// if err != nil {
	// 	return fmt.Errorf("checking email existence failed: %w", err)
	// }
	// if exists {
	// 	return fmt.Errorf("email %s already in use", user.Email)
	// }

	return s.userRepo.Create(ctx, user)
}

func (s *userService) Get(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	if !user.Status.IsValid() {
		return fmt.Errorf("invalid status %q", user.Status)
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
