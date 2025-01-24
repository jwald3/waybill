package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jwald3/waybill/internal/database"
	"github.com/jwald3/waybill/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type incidentReportRepository struct {
	incidentReports *mongo.Collection
}

type IncidentReportRepository interface {
	Create(ctx context.Context, incidentReport *domain.IncidentReport) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.IncidentReport, error)
	Update(ctx context.Context, incidentReport *domain.IncidentReport) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.IncidentReport, error)
}

func NewIncidentReportRepository(db *database.MongoDB) IncidentReportRepository {
	return &incidentReportRepository{
		incidentReports: db.Database.Collection("incident_reports"),
	}
}

func (r *incidentReportRepository) Create(ctx context.Context, incidentReport *domain.IncidentReport) error {
	now := time.Now()
	incidentReport.CreatedAt = primitive.NewDateTimeFromTime(now)
	incidentReport.UpdatedAt = primitive.NewDateTimeFromTime(now)

	_, err := r.incidentReports.InsertOne(ctx, incidentReport)
	if err != nil {
		return fmt.Errorf("failed to create incidentReport: %w", err)
	}

	return nil
}

func (r *incidentReportRepository) GetById(ctx context.Context, id primitive.ObjectID) (*domain.IncidentReport, error) {
	filter := bson.M{"_id": id}

	var incidentReport domain.IncidentReport
	err := r.incidentReports.FindOne(ctx, filter).Decode(&incidentReport)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get incidentReport: %w", err)
	}
	return &incidentReport, nil
}

func (r *incidentReportRepository) Update(ctx context.Context, incidentReport *domain.IncidentReport) error {
	filter := bson.M{"_id": incidentReport.ID}
	update := bson.M{
		"$set": bson.M{
			"truck_id":        incidentReport.TruckID,
			"driver_id":       incidentReport.DriverID,
			"type":            incidentReport.Type,
			"description":     incidentReport.Description,
			"date":            incidentReport.Date,
			"location":        incidentReport.Location,
			"damage_estimate": incidentReport.DamageEstimate,
			"updated_at":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err := r.incidentReports.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update incidentReport: %w", err)
	}

	return nil
}

func (r *incidentReportRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.incidentReports.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete incidentReport: %w", err)
	}
	return nil
}

func (r *incidentReportRepository) List(ctx context.Context, limit, offset int64) ([]*domain.IncidentReport, error) {
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := r.incidentReports.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve list of users: %w", err)
	}
	defer cursor.Close(ctx)

	var incidentReports []*domain.IncidentReport
	for cursor.Next(ctx) {
		var d domain.IncidentReport
		if err := cursor.Decode(&d); err != nil {
			return nil, fmt.Errorf("failed to decode incidentReport: %w", err)
		}
		incidentReports = append(incidentReports, &d)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return incidentReports, nil
}
