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
)

type incidentReportRepository struct {
	incidentReports *mongo.Collection
}

type IncidentReportRepository interface {
	Create(ctx context.Context, incidentReport *domain.IncidentReport) error
	GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.IncidentReport, error)
	Update(ctx context.Context, incidentReport *domain.IncidentReport) error
	Delete(ctx context.Context, id, userID primitive.ObjectID) error
	List(ctx context.Context, filter domain.IncidentReportFilter) (*ListIncidentReportsResult, error)
}

type ListIncidentReportsResult struct {
	IncidentReports []*domain.IncidentReport
	Total           int64
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

func (r *incidentReportRepository) GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.IncidentReport, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": id, "user_id": userID}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "trips",
			"localField":   "trip_id",
			"foreignField": "_id",
			"as":           "trip",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$trip",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "trucks",
			"localField":   "truck_id",
			"foreignField": "_id",
			"as":           "truck",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$truck",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "drivers",
			"localField":   "driver_id",
			"foreignField": "_id",
			"as":           "driver",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$driver",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"trip_id":   0,
			"truck_id":  0,
			"driver_id": 0,
		}}},
	}

	var result domain.IncidentReport
	cursor, err := r.incidentReports.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		if cursor.Err() != nil {
			return nil, fmt.Errorf("cursor error: %w", cursor.Err())
		}
		return nil, nil
	}

	if err := cursor.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode driver: %w", err)
	}

	return &result, nil
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

	result, err := r.incidentReports.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update incidentReport: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("incident report not found")
	}

	return nil
}

func (r *incidentReportRepository) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	result, err := r.incidentReports.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		return fmt.Errorf("failed to delete incidentReport: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrIncidentReportNotFound
	}

	return nil
}

func (r *incidentReportRepository) List(ctx context.Context, filter domain.IncidentReportFilter) (*ListIncidentReportsResult, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	filterQuery := bson.M{"user_id": filter.UserID}

	if filter.TripID != nil {
		filterQuery["trip_id"] = filter.TripID
	}

	if filter.TruckID != nil {
		filterQuery["truck_id"] = filter.TruckID
	}

	if filter.DriverID != nil {
		filterQuery["driver_id"] = filter.DriverID
	}

	if filter.Type != "" {
		filterQuery["type"] = filter.Type
	}

	total, err := r.incidentReports.CountDocuments(ctx, filterQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filterQuery}},
		{{Key: "$sort", Value: bson.M{"_id": -1}}},
		{{Key: "$skip", Value: filter.Offset}},
		{{Key: "$limit", Value: filter.Limit}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "trips",
			"localField":   "trip_id",
			"foreignField": "_id",
			"as":           "trip",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$trip",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "trucks",
			"localField":   "truck_id",
			"foreignField": "_id",
			"as":           "truck",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$truck",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "drivers",
			"localField":   "driver_id",
			"foreignField": "_id",
			"as":           "driver",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$driver",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"trip_id":   0,
			"truck_id":  0,
			"driver_id": 0,
		}}},
	}

	cursor, err := r.incidentReports.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate query: %w", err)
	}
	defer cursor.Close(ctx)

	incidentReports := make([]*domain.IncidentReport, 0, filter.Limit)
	if err := cursor.All(ctx, &incidentReports); err != nil {
		return nil, fmt.Errorf("failed to decode incident reports: %w", err)
	}

	return &ListIncidentReportsResult{
		IncidentReports: incidentReports,
		Total:           total,
	}, nil
}
