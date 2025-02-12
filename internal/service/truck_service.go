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
	truckNotFound = "unable to retrieve truck: %w"
)

type TruckService interface {
	Create(ctx context.Context, truck *domain.Truck) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.Truck, error)
	Update(ctx context.Context, truck *domain.Truck) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) (*repository.ListTrucksResult, error)
	UpdateTruckStatus(ctx context.Context, id primitive.ObjectID, newStatus domain.TruckStatus) error
	UpdateTruckMileage(ctx context.Context, id primitive.ObjectID, newMileage int) error
	UpdateTruckMaintenance(ctx context.Context, id primitive.ObjectID, lastMaintenance string) error
}

type truckService struct {
	db        *database.MongoDB
	truckRepo repository.TruckRepository
}

func NewTruckService(db *database.MongoDB, truckRepo repository.TruckRepository) TruckService {
	return &truckService{
		db:        db,
		truckRepo: truckRepo,
	}
}

func (s *truckService) Create(ctx context.Context, truck *domain.Truck) error {
	if err := s.truckRepo.Create(ctx, truck); err != nil {
		return fmt.Errorf("failed to create truck: %w", err)
	}

	return nil
}

func (s *truckService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Truck, error) {
	truck, err := s.truckRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(truckNotFound, err)
	}
	if truck == nil {
		return nil, fmt.Errorf("truck with ID %v not found", id)
	}

	return truck, nil
}

func (s *truckService) Update(ctx context.Context, truck *domain.Truck) error {
	err := s.truckRepo.Update(ctx, truck)
	if err != nil {
		return fmt.Errorf(truckNotFound, err)
	}

	return nil
}

func (s *truckService) Delete(ctx context.Context, id primitive.ObjectID) error {
	if err := s.truckRepo.Delete(ctx, id); err != nil {
		if err == domain.ErrTruckNotFound {
			return err
		}

		return fmt.Errorf("failed to delete truck: %w", err)
	}

	return nil
}

func (s *truckService) List(ctx context.Context, limit, offset int64) (*repository.ListTrucksResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.truckRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list trucks: %w", err)
	}

	if result.Trucks == nil {
		result.Trucks = []*domain.Truck{}
	}

	return result, nil
}

// atomic methods

func (s *truckService) UpdateTruckStatus(ctx context.Context, id primitive.ObjectID, newStatus domain.TruckStatus) error {
	truck, err := s.truckRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf(truckNotFound, err)
	}

	if err := truck.Status.CanTransitionTo(newStatus); err != nil {
		return err
	}

	truck.Status = newStatus

	return s.truckRepo.Update(ctx, truck)
}

func (s *truckService) UpdateTruckMileage(ctx context.Context, id primitive.ObjectID, newMileage int) error {
	truck, err := s.truckRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf(truckNotFound, err)
	}

	truck.Mileage = newMileage

	return s.truckRepo.Update(ctx, truck)
}

func (s *truckService) UpdateTruckMaintenance(ctx context.Context, id primitive.ObjectID, lastMaintenance string) error {
	truck, err := s.truckRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf(truckNotFound, err)
	}

	truck.LastMaintenance = lastMaintenance

	return s.truckRepo.Update(ctx, truck)
}
