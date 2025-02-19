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

type maintenanceLogRepository struct {
	maintenanceLogs *mongo.Collection
}

type MaintenanceLogRepository interface {
	Create(ctx context.Context, maintenanceLog *domain.MaintenanceLog) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.MaintenanceLog, error)
	Update(ctx context.Context, maintenanceLog *domain.MaintenanceLog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter domain.MaintenanceLogFilter) (*ListMaintenanceLogsResult, error)
}

type ListMaintenanceLogsResult struct {
	MaintenanceLogs []*domain.MaintenanceLog
	Total           int64
}

func NewMaintenanceLogRepository(db *database.MongoDB) MaintenanceLogRepository {
	return &maintenanceLogRepository{
		maintenanceLogs: db.Database.Collection("maintenance_logs"),
	}
}

func (r *maintenanceLogRepository) Create(ctx context.Context, maintenanceLog *domain.MaintenanceLog) error {
	now := time.Now()
	maintenanceLog.CreatedAt = primitive.NewDateTimeFromTime(now)
	maintenanceLog.UpdatedAt = primitive.NewDateTimeFromTime(now)

	_, err := r.maintenanceLogs.InsertOne(ctx, maintenanceLog)
	if err != nil {
		return fmt.Errorf("failed to create maintenance log %w", err)
	}

	return nil
}

func (r *maintenanceLogRepository) GetById(ctx context.Context, id primitive.ObjectID) (*domain.MaintenanceLog, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": id}}},
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
		{{Key: "$project", Value: bson.M{
			"truck_id": 0,
		}}},
	}

	var result domain.MaintenanceLog
	cursor, err := r.maintenanceLogs.Aggregate(ctx, pipeline)
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
		return nil, fmt.Errorf("failed to decode maintenance log: %w", err)
	}

	return &result, nil
}

func (r *maintenanceLogRepository) Update(ctx context.Context, maintenanceLog *domain.MaintenanceLog) error {
	filter := bson.M{"_id": maintenanceLog.ID}
	update := bson.M{
		"$set": bson.M{
			"truck_id":     maintenanceLog.TruckID,
			"date":         maintenanceLog.Date,
			"service_type": maintenanceLog.ServiceType,
			"cost":         maintenanceLog.Cost,
			"notes":        maintenanceLog.Notes,
			"mechanic":     maintenanceLog.Mechanic,
			"location":     maintenanceLog.Location,
			"updated_at":   primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.maintenanceLogs.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update maintenance log %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("maintenance log not found")
	}

	return nil
}

func (r *maintenanceLogRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.maintenanceLogs.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete maintenance log %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrMaintenanceLogNotFound
	}

	return nil
}

func (r *maintenanceLogRepository) List(ctx context.Context, filter domain.MaintenanceLogFilter) (*ListMaintenanceLogsResult, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	filterQuery := bson.M{}

	if filter.TruckID != &primitive.NilObjectID {
		filterQuery["truck_id"] = filter.TruckID
	}

	if filter.ServiceType != "" {
		filterQuery["service_type"] = filter.ServiceType
	}

	total, err := r.maintenanceLogs.CountDocuments(ctx, filterQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.M{"_id": -1}}},
		{{Key: "$skip", Value: filter.Offset}},
		{{Key: "$limit", Value: filter.Limit}},
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
		{{Key: "$project", Value: bson.M{
			"truck_id": 0,
		}}},
	}

	cursor, err := r.maintenanceLogs.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate query: %w", err)
	}
	defer cursor.Close(ctx)

	maintenanceLogs := make([]*domain.MaintenanceLog, 0, filter.Limit)
	if err := cursor.All(ctx, &maintenanceLogs); err != nil {
		return nil, fmt.Errorf("failed to decode maintenance logs: %w", err)
	}

	return &ListMaintenanceLogsResult{
		MaintenanceLogs: maintenanceLogs,
		Total:           total,
	}, nil
}
