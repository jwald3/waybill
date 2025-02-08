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
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.Driver, error)
	Update(ctx context.Context, driver *domain.Driver) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) (*repository.ListDriversResult, error)
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
	if !driver.EmploymentStatus.IsValid() {
		return fmt.Errorf("invalid employment status: %s", driver.EmploymentStatus)
	}

	if err := s.driverRepo.Create(ctx, driver); err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	return nil
}

func (s *driverService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Driver, error) {
	driver, err := s.driverRepo.GetById(ctx, id)
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

func (s *driverService) Delete(ctx context.Context, id primitive.ObjectID) error {
	if err := s.driverRepo.Delete(ctx, id); err != nil {
		if err == domain.ErrDriverNotFound {
			return err
		}
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	return nil
}

func (s *driverService) List(ctx context.Context, limit, offset int64) (*repository.ListDriversResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.driverRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list drivers: %w", err)
	}

	if result.Drivers == nil {
		result.Drivers = []*domain.Driver{}
	}

	return result, nil
}

// Atomic methods
func (s *driverService) UpdateEmploymentStatus(ctx context.Context, id primitive.ObjectID, newStatus domain.EmploymentStatus) error {
	driver, err := s.driverRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf(driverNotFound, err)
	}
	if err := driver.ChangeEmploymentStatus(newStatus); err != nil {
		return err
	}
	return s.driverRepo.Update(ctx, driver)
}
