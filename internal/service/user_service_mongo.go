package service

import (
	"context"
	"fmt"

	"github.com/jwald3/go_rest_template/internal/database"
	"github.com/jwald3/go_rest_template/internal/domain"
	"github.com/jwald3/go_rest_template/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userMongoService struct {
	mdb      *database.MongoDB
	userRepo repository.UserMongoRepository
}

type UserMongoService interface {
	Create(ctx context.Context, user *domain.User) error
	Get(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.User, error)
}

func NewUserMongoService(mdb *database.MongoDB, userRepo repository.UserMongoRepository) UserMongoService {
	return &userMongoService{
		mdb:      mdb,
		userRepo: userRepo,
	}
}

func (s *userMongoService) Create(ctx context.Context, user *domain.User) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("email duplicate check failed: %w", err)
	}

	if exists {
		return fmt.Errorf("email %s is already in use", user.Email)
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *userMongoService) Get(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}
	return user, nil
}

func (s *userMongoService) Update(ctx context.Context, user *domain.User) error {
	if !user.Status.IsValid() {
		return fmt.Errorf("invalid status %q", user.Status)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *userMongoService) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.userRepo.GetByID(ctx, id)

	if err != nil {
		return fmt.Errorf("failed to retrieve user: %w", err)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user")
	}

	return nil
}

func (s *userMongoService) List(ctx context.Context, limit, offset int64) ([]*domain.User, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
