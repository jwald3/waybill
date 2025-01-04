package service

import (
	"context"
	"fmt"

	"github.com/jwald3/go_rest_template/internal/database"
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
	db       *database.DB
	userRepo repository.UserRepository
}

func NewUserService(db *database.DB, userRepo repository.UserRepository) UserService {
	return &userService{
		// we include the database in the userService so we have access to the transaction methods.
		// if you do not need transaction support, do not worry about including this
		db:       db,
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, user *domain.User) error {
	// some routes' logic will require a series of actions to occur within a transaction in an all-or-nothing type of payload.
	// either all items within the transaction pass and are saved or none are saved and any changes are rolled back. This type of
	// architecture makes the most sense when you're including auditing or changing multiple resources in one sequence.
	return s.db.ExecuteTx(ctx, func(tx *database.Transaction) error {

		exists, err := s.userRepo.ExistsByEmailTx(ctx, tx, user.Email)

		if err != nil {
			return fmt.Errorf("checking email existence failed: %w", err)
		}

		if exists {
			return fmt.Errorf("email %s already in use", user.Email)
		}

		// we use the "transaction-friendly" method to ensure the action occurs along the same strand as all other transaction members
		if err := s.userRepo.CreateTx(ctx, tx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// here is where you could add additional transaction members, such as updating an audit log, performing external services, etc.

		return nil
	})

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
	return s.db.ExecuteTx(ctx, func(tx *database.Transaction) error {
		if !user.Status.IsValid() {
			return fmt.Errorf("invalid status %q", user.Status)
		}

		if err := s.userRepo.UpdateTx(ctx, tx, user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		// additional transaction members here...

		return nil
	})
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	return s.db.ExecuteTx(ctx, func(tx *database.Transaction) error {
		_, err := s.userRepo.GetByID(ctx, id)

		if err != nil {
			return fmt.Errorf("failed to retrieve user: %w", err)
		}

		if err := s.userRepo.DeleteTx(ctx, tx, id); err != nil {
			return fmt.Errorf("failed to delete user")
		}

		// additional transaction members here...

		return nil
	})
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
