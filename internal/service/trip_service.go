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
	tripNotFound = "unable to retrieve trip: %w"
)

type TripService interface {
	Create(ctx context.Context, trip *domain.Trip) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.Trip, error)
	Update(ctx context.Context, trip *domain.Trip) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) (*repository.ListTripsResult, error)
}

type tripService struct {
	db       *database.MongoDB
	tripRepo repository.TripRepository
}

func NewTripService(db *database.MongoDB, tripRepo repository.TripRepository) TripService {
	return &tripService{
		db:       db,
		tripRepo: tripRepo,
	}
}

func (s *tripService) Create(ctx context.Context, trip *domain.Trip) error {
	if err := s.tripRepo.Create(ctx, trip); err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}

	return nil
}

func (s *tripService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Trip, error) {
	trip, err := s.tripRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(tripNotFound, err)
	}
	if trip == nil {
		return nil, fmt.Errorf("trip with ID %v not found", id)
	}

	return trip, nil
}

func (s *tripService) Update(ctx context.Context, trip *domain.Trip) error {
	err := s.tripRepo.Update(ctx, trip)
	if err != nil {
		return fmt.Errorf(tripNotFound, err)
	}

	return nil
}

func (s *tripService) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.tripRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf(tripNotFound, err)
	}

	if err := s.tripRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	return nil
}

func (s *tripService) List(ctx context.Context, limit, offset int64) (*repository.ListTripsResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.tripRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list trips: %w", err)
	}

	if result.Trips == nil {
		result.Trips = []*domain.Trip{}
	}

	return result, nil
}
