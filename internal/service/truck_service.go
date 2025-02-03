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
	List(ctx context.Context, limit, offset int64) ([]*domain.Truck, error)
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
	_, err := s.truckRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf(truckNotFound, err)
	}

	if err := s.truckRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete truck: %w", err)
	}

	return nil
}

func (s *truckService) List(ctx context.Context, limit, offset int64) ([]*domain.Truck, error) {
	trucks, err := s.truckRepo.List(ctx, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to list trucks: %w", err)
	}

	return trucks, nil
}

// atomic methods

func (s *truckService) UpdateEmploymentStatus(ctx context.Context, id primitive.ObjectID, newStatus domain.TruckStatus) error {
	truck, err := s.truckRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf(truckNotFound, err)
	}
	if err := truck.ChangeTruckStatus(newStatus); err != nil {
		return err
	}
	return s.truckRepo.Update(ctx, truck)
}
