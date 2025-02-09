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
	incidentReportNotFound = "unable to retrieve incident report: %w"
)

type IncidentReportService interface {
	Create(ctx context.Context, incidentReport *domain.IncidentReport) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.IncidentReport, error)
	Update(ctx context.Context, incidentReport *domain.IncidentReport) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) (*repository.ListIncidentReportsResult, error)
}

type incidentReportService struct {
	db                 *database.MongoDB
	incidentReportRepo repository.IncidentReportRepository
}

func NewIncidentReportService(db *database.MongoDB, incidentReportRepo repository.IncidentReportRepository) IncidentReportService {
	return &incidentReportService{
		db:                 db,
		incidentReportRepo: incidentReportRepo,
	}
}

func (s *incidentReportService) Create(ctx context.Context, incidentReport *domain.IncidentReport) error {
	if err := s.incidentReportRepo.Create(ctx, incidentReport); err != nil {
		return fmt.Errorf("failed to create incident report: %w", err)
	}

	return nil
}

func (s *incidentReportService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.IncidentReport, error) {
	incidentReport, err := s.incidentReportRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(incidentReportNotFound, err)
	}
	if incidentReport == nil {
		return nil, fmt.Errorf("incident report with ID %v not found", id)
	}

	return incidentReport, nil
}

func (s *incidentReportService) Update(ctx context.Context, incidentReport *domain.IncidentReport) error {
	err := s.incidentReportRepo.Update(ctx, incidentReport)
	if err != nil {
		return fmt.Errorf(incidentReportNotFound, err)
	}

	return nil
}

func (s *incidentReportService) Delete(ctx context.Context, id primitive.ObjectID) error {
	if err := s.incidentReportRepo.Delete(ctx, id); err != nil {
		if err == domain.ErrIncidentReportNotFound {
			return err
		}
		return fmt.Errorf("failed to delete incident report: %w", err)
	}

	return nil
}

func (s *incidentReportService) List(ctx context.Context, limit, offset int64) (*repository.ListIncidentReportsResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	result, err := s.incidentReportRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list incident reports: %w", err)
	}

	if result.IncidentReports == nil {
		result.IncidentReports = []*domain.IncidentReport{}
	}

	return result, nil
}
