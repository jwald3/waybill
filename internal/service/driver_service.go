package service

import (
	"context"
	"fmt"

	"github.com/jwald3/waybill/internal/database"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	driverNotFound = "unable to retrieve driver: %w"
)

type DriverService interface {
	Create(ctx context.Context, driver *domain.Driver) error
	GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.Driver, error)
	Update(ctx context.Context, driver *domain.Driver) error
	Delete(ctx context.Context, id, userID primitive.ObjectID) error
	List(ctx context.Context, filter domain.DriverFilter) (*repository.ListDriversResult, error)
	SuspendDriver(ctx context.Context, id, userID primitive.ObjectID) error
	TerminateDriver(ctx context.Context, id, userID primitive.ObjectID) error
	ActivateDriver(ctx context.Context, id, userID primitive.ObjectID) error
}

type driverService struct {
	db         *database.MongoDB
	driverRepo repository.DriverRepository
}

func NewDriverService(db *database.MongoDB, driverRepo repository.DriverRepository) DriverService {
	return &driverService{
		db:         db,
		driverRepo: driverRepo,
	}
}

func (s *driverService) Create(ctx context.Context, driver *domain.Driver) error {
	if err := s.driverRepo.Create(ctx, driver); err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	return nil
}

func (s *driverService) GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.Driver, error) {
	driver, err := s.driverRepo.GetById(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf(driverNotFound, err)
	}
	if driver == nil {
		return nil, fmt.Errorf("driver with ID %v not found", id)
	}

	return driver, nil
}

func (s *driverService) Update(ctx context.Context, driver *domain.Driver) error {
	err := s.driverRepo.Update(ctx, driver)
	if err != nil {
		return fmt.Errorf(driverNotFound, err)
	}

	return nil
}

func (s *driverService) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	if err := s.driverRepo.Delete(ctx, id, userID); err != nil {
		if err == domain.ErrDriverNotFound {
			return err
		}
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	return nil
}

func (s *driverService) List(ctx context.Context, filter domain.DriverFilter) (*repository.ListDriversResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.driverRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list drivers: %w", err)
	}

	if result.Drivers == nil {
		result.Drivers = []*domain.Driver{}
	}

	return result, nil
}

// Atomic methods
func (s *driverService) SuspendDriver(ctx context.Context, id, userID primitive.ObjectID) error {
	driver, err := s.driverRepo.GetById(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to get driver: %w", err)
	}
	if driver == nil {
		return domain.ErrDriverNotFound
	}

	if err := driver.InitializeStateMachine(); err != nil {
		return fmt.Errorf("failed to initialize state machine: %w", err)
	}

	if err := driver.SuspendDriver(); err != nil {
		return err
	}

	return s.driverRepo.Update(ctx, driver)
}

func (s *driverService) TerminateDriver(ctx context.Context, id, userID primitive.ObjectID) error {
	driver, err := s.driverRepo.GetById(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to get driver: %w", err)
	}
	if driver == nil {
		return domain.ErrDriverNotFound
	}

	if err := driver.InitializeStateMachine(); err != nil {
		return fmt.Errorf("failed to initialize state machine: %w", err)
	}

	if err := driver.TerminateDriver(); err != nil {
		return err
	}

	return s.driverRepo.Update(ctx, driver)
}

func (s *driverService) ActivateDriver(ctx context.Context, id, userID primitive.ObjectID) error {
	driver, err := s.driverRepo.GetById(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to get driver: %w", err)
	}
	if driver == nil {
		return domain.ErrDriverNotFound
	}

	if err := driver.InitializeStateMachine(); err != nil {
		return fmt.Errorf("failed to initialize state machine: %w", err)
	}

	if err := driver.ActivateDriver(); err != nil {
		return err
	}

	return s.driverRepo.Update(ctx, driver)
}
