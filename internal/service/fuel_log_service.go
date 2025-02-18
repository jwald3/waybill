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
	fuelLogNotFound = "unable to retrieve fuel log: %w"
)

type FuelLogService interface {
	Create(ctx context.Context, fuelLog *domain.FuelLog) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.FuelLog, error)
	Update(ctx context.Context, fuelLog *domain.FuelLog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter domain.FuelLogFilter) (*repository.ListFuelLogsResult, error)
}

type fuelLogService struct {
	db          *database.MongoDB
	fuelLogRepo repository.FuelLogRepository
}

func NewFuelLogService(db *database.MongoDB, fuelLogRepo repository.FuelLogRepository) FuelLogService {
	return &fuelLogService{
		db:          db,
		fuelLogRepo: fuelLogRepo,
	}
}

func (s *fuelLogService) Create(ctx context.Context, fuelLog *domain.FuelLog) error {
	if err := s.fuelLogRepo.Create(ctx, fuelLog); err != nil {
		return fmt.Errorf("failed to create fuel log: %w", err)
	}

	return nil
}

func (s *fuelLogService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.FuelLog, error) {
	fuelLog, err := s.fuelLogRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(fuelLogNotFound, err)
	}
	if fuelLog == nil {
		return nil, fmt.Errorf("fuel log with ID %v not found", id)
	}

	return fuelLog, nil
}

func (s *fuelLogService) Update(ctx context.Context, fuelLog *domain.FuelLog) error {
	err := s.fuelLogRepo.Update(ctx, fuelLog)
	if err != nil {
		return fmt.Errorf(fuelLogNotFound, err)
	}

	return nil
}

func (s *fuelLogService) Delete(ctx context.Context, id primitive.ObjectID) error {
	if err := s.fuelLogRepo.Delete(ctx, id); err != nil {
		if err == domain.ErrFuelLogNotFound {
			return err
		}
		return fmt.Errorf("failed to delete fuel log: %w", err)
	}

	return nil
}

func (s *fuelLogService) List(ctx context.Context, filter domain.FuelLogFilter) (*repository.ListFuelLogsResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.fuelLogRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list fuel logs: %w", err)
	}

	if result.FuelLogs == nil {
		result.FuelLogs = []*domain.FuelLog{}
	}

	return result, nil
}
