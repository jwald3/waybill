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
	maintenanceLogNotFound = "unable to retrieve maintenance log: %w"
)

type MaintenanceLogService interface {
	Create(ctx context.Context, maintenanceLog *domain.MaintenanceLog) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.MaintenanceLog, error)
	Update(ctx context.Context, maintenanceLog *domain.MaintenanceLog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter domain.MaintenanceLogFilter) (*repository.ListMaintenanceLogsResult, error)
}

type maintenanceLogService struct {
	db                 *database.MongoDB
	maintenanceLogRepo repository.MaintenanceLogRepository
}

func NewMaintenanceLogService(db *database.MongoDB, maintenanceLogRepo repository.MaintenanceLogRepository) MaintenanceLogService {
	return &maintenanceLogService{
		db:                 db,
		maintenanceLogRepo: maintenanceLogRepo,
	}
}

func (s *maintenanceLogService) Create(ctx context.Context, maintenanceLog *domain.MaintenanceLog) error {
	if err := s.maintenanceLogRepo.Create(ctx, maintenanceLog); err != nil {
		return fmt.Errorf("failed to create maintenance log: %w", err)
	}

	return nil
}

func (s *maintenanceLogService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.MaintenanceLog, error) {
	maintenanceLog, err := s.maintenanceLogRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(maintenanceLogNotFound, err)
	}
	if maintenanceLog == nil {
		return nil, fmt.Errorf("maintenance log with ID %v not found", id)
	}

	return maintenanceLog, nil
}

func (s *maintenanceLogService) Update(ctx context.Context, maintenanceLog *domain.MaintenanceLog) error {
	err := s.maintenanceLogRepo.Update(ctx, maintenanceLog)
	if err != nil {
		return fmt.Errorf(maintenanceLogNotFound, err)
	}

	return nil
}

func (s *maintenanceLogService) Delete(ctx context.Context, id primitive.ObjectID) error {
	if err := s.maintenanceLogRepo.Delete(ctx, id); err != nil {
		if err == domain.ErrMaintenanceLogNotFound {
			return err
		}

		return fmt.Errorf("failed to delete maintenance log: %w", err)
	}

	return nil
}

func (s *maintenanceLogService) List(ctx context.Context, filter domain.MaintenanceLogFilter) (*repository.ListMaintenanceLogsResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.maintenanceLogRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list maintenance logs: %w", err)
	}

	if result.MaintenanceLogs == nil {
		result.MaintenanceLogs = []*domain.MaintenanceLog{}
	}

	return result, nil
}
