package service

import (
	"context"
	"fmt"
	"time"

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
	GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.Trip, error)
	Update(ctx context.Context, trip *domain.Trip) error
	Delete(ctx context.Context, id, userID primitive.ObjectID) error
	List(ctx context.Context, filter domain.TripFilter) (*repository.ListTripsResult, error)
	AddNote(ctx context.Context, id, userID primitive.ObjectID, content string) error
	BeginTrip(ctx context.Context, id, userID primitive.ObjectID, departureTime time.Time) error
	CancelTrip(ctx context.Context, id, userID primitive.ObjectID) error
	FinishTripSuccessfully(ctx context.Context, id, userID primitive.ObjectID, arrivalTime time.Time) error
	FinishTripUnsuccessfully(ctx context.Context, id, userID primitive.ObjectID, arrivalTime time.Time) error
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

func (s *tripService) GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.Trip, error) {
	trip, err := s.tripRepo.GetById(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf(tripNotFound, err)
	}
	if trip == nil {
		return nil, fmt.Errorf("trip with ID %v not found", id)
	}

	if trip.Notes == nil {
		trip.Notes = make([]domain.TripNote, 0)
	}

	// Initialize the state machine for the retrieved trip
	if err := trip.InitializeStateMachine(); err != nil {
		return nil, fmt.Errorf("failed to initialize state machine: %w", err)
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

func (s *tripService) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	if err := s.tripRepo.Delete(ctx, id, userID); err != nil {
		if err == domain.ErrTripNotFound {
			return err
		}

		return fmt.Errorf("failed to delete trip: %w", err)
	}

	return nil
}

func (s *tripService) List(ctx context.Context, filter domain.TripFilter) (*repository.ListTripsResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.tripRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list trips: %w", err)
	}

	if result.Trips == nil {
		result.Trips = []*domain.Trip{}
	}

	return result, nil
}

func (s *tripService) AddNote(ctx context.Context, id, userID primitive.ObjectID, content string) error {
	trip, err := s.tripRepo.GetById(ctx, id, userID)

	if trip == nil {
		return fmt.Errorf("trip with ID %v not found", id)
	}

	if err != nil {
		return fmt.Errorf(tripNotFound, err)
	}

	if err := trip.AddNote(content); err != nil {
		return err
	}

	return s.tripRepo.Update(ctx, trip)
}

func (s *tripService) CancelTrip(ctx context.Context, id, userID primitive.ObjectID) error {
	trip, err := s.tripRepo.GetById(ctx, id, userID)
	if err != nil {
		return fmt.Errorf(tripNotFound, err)
	}

	if trip == nil {
		return fmt.Errorf("trip with ID %v not found", id)
	}

	if err := trip.InitializeStateMachine(); err != nil {
		return fmt.Errorf("failed to initialize state machine: %w", err)
	}

	if err := trip.CancelTrip(); err != nil {
		return err
	}

	return s.tripRepo.Update(ctx, trip)
}

func (s *tripService) BeginTrip(ctx context.Context, id, userID primitive.ObjectID, departureTime time.Time) error {
	trip, err := s.tripRepo.GetById(ctx, id, userID)
	if err != nil {
		return fmt.Errorf(tripNotFound, err)
	}

	if trip == nil {
		return fmt.Errorf("trip with ID %v not found", id)
	}

	// Initialize the state machine for the retrieved trip
	if err := trip.InitializeStateMachine(); err != nil {
		return fmt.Errorf("failed to initialize state machine: %w", err)
	}

	if err := trip.BeginTrip(departureTime); err != nil {
		return fmt.Errorf("an error occurred when attempting to begin trip: %w", err)
	}

	// Update with the full trip object that contains all references
	return s.tripRepo.Update(ctx, trip)
}

func (s *tripService) FinishTripSuccessfully(ctx context.Context, id, userID primitive.ObjectID, arrivalTime time.Time) error {
	trip, err := s.tripRepo.GetById(ctx, id, userID)
	if err != nil {
		return fmt.Errorf(tripNotFound, err)
	}

	if trip == nil {
		return fmt.Errorf("trip with ID %v not found", id)
	}

	if err := trip.InitializeStateMachine(); err != nil {
		return fmt.Errorf("failed to initialize state machine: %w", err)
	}

	if err := trip.CompleteTripSuccessfully(arrivalTime); err != nil {
		return fmt.Errorf("an error occurred when attempting to complete trip: %w", err)
	}

	return s.tripRepo.Update(ctx, trip)
}

func (s *tripService) FinishTripUnsuccessfully(ctx context.Context, id, userID primitive.ObjectID, arrivalTime time.Time) error {
	trip, err := s.tripRepo.GetById(ctx, id, userID)
	if err != nil {
		return fmt.Errorf(tripNotFound, err)
	}

	if trip == nil {
		return fmt.Errorf("trip with ID %v not found", id)
	}

	if err := trip.InitializeStateMachine(); err != nil {
		return fmt.Errorf("failed to initialize state machine: %w", err)
	}

	if err := trip.CompleteTripUnsuccessfully(arrivalTime); err != nil {
		return fmt.Errorf("an error occurred when attempting to complete trip: %w", err)
	}

	return s.tripRepo.Update(ctx, trip)
}
