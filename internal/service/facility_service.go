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
	facilityNotFound = "unable to retrieve facility: %w"
)

type FacilityService interface {
	Create(ctx context.Context, facility *domain.Facility) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.Facility, error)
	Update(ctx context.Context, facility *domain.Facility) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) (*repository.ListFacilitiesResult, error)
	UpdateAvailableFacilityServices(ctx context.Context, id primitive.ObjectID, servicesAvailable []domain.FacilityService) error
}

type facilityService struct {
	db           *database.MongoDB
	facilityRepo repository.FacilityRepository
}

func NewFacilityService(db *database.MongoDB, facilityRepo repository.FacilityRepository) FacilityService {
	return &facilityService{
		db:           db,
		facilityRepo: facilityRepo,
	}
}

func (s *facilityService) Create(ctx context.Context, facility *domain.Facility) error {
	if err := s.facilityRepo.Create(ctx, facility); err != nil {
		return fmt.Errorf("failed to create facility: %w", err)
	}

	return nil
}

func (s *facilityService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Facility, error) {
	facility, err := s.facilityRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(facilityNotFound, err)
	}
	if facility == nil {
		return nil, fmt.Errorf("facility with ID %v not found", id)
	}

	return facility, nil
}

func (s *facilityService) Update(ctx context.Context, facility *domain.Facility) error {
	err := s.facilityRepo.Update(ctx, facility)
	if err != nil {
		return fmt.Errorf(facilityNotFound, err)
	}

	return nil
}

func (s *facilityService) Delete(ctx context.Context, id primitive.ObjectID) error {
	if err := s.facilityRepo.Delete(ctx, id); err != nil {
		if err == domain.ErrFacilityNotFound {
			return err
		}
		return fmt.Errorf("failed to delete facility: %w", err)
	}

	return nil
}

func (s *facilityService) List(ctx context.Context, limit, offset int64) (*repository.ListFacilitiesResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.facilityRepo.List(ctx, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to list facilities: %w", err)
	}

	if result.Facilities == nil {
		result.Facilities = []*domain.Facility{}
	}

	return result, nil
}

// atomic methods
func (s *facilityService) UpdateAvailableFacilityServices(ctx context.Context, id primitive.ObjectID, servicesAvailable []domain.FacilityService) error {
	// Validate services first
	for _, service := range servicesAvailable {
		if !service.IsValid() {
			return fmt.Errorf("invalid facility service: %s", service)
		}
	}

	// Check if facility exists
	facility, err := s.facilityRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get facility: %w", err)
	}
	if facility == nil {
		return domain.ErrFacilityNotFound
	}

	// Update the services
	if err := s.facilityRepo.UpdateAvailableFacilityServices(ctx, id, servicesAvailable); err != nil {
		return fmt.Errorf("failed to update available facility services: %w", err)
	}

	return nil
}
