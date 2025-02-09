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

type fuelLogRepository struct {
	fuelLogs *mongo.Collection
}

type FuelLogRepository interface {
	Create(ctx context.Context, fuelLog *domain.FuelLog) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.FuelLog, error)
	Update(ctx context.Context, fuelLog *domain.FuelLog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) (*ListFuelLogsResult, error)
}

type ListFuelLogsResult struct {
	FuelLogs []*domain.FuelLog
	Total    int64
}

func NewFuelLogRepository(db *database.MongoDB) FuelLogRepository {
	return &fuelLogRepository{
		fuelLogs: db.Database.Collection("fuel_logs"),
	}
}

func (r *fuelLogRepository) Create(ctx context.Context, fuelLog *domain.FuelLog) error {
	now := time.Now()
	fuelLog.CreatedAt = primitive.NewDateTimeFromTime(now)
	fuelLog.UpdatedAt = primitive.NewDateTimeFromTime(now)

	_, err := r.fuelLogs.InsertOne(ctx, fuelLog)
	if err != nil {
		return fmt.Errorf("failed to create fuelLog: %w", err)
	}

	return nil
}

func (r *fuelLogRepository) GetById(ctx context.Context, id primitive.ObjectID) (*domain.FuelLog, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"_id": id,
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
		{{
			Key: "$project", Value: bson.M{
				"truck_id":  0,
				"driver_id": 0,
			},
		}},
	}

	var result domain.FuelLog
	cursor, err := r.fuelLogs.Aggregate(ctx, pipeline)
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
		return nil, fmt.Errorf("failed to decode fuel log: %w", err)
	}

	return &result, nil
}

func (r *fuelLogRepository) Update(ctx context.Context, fuelLog *domain.FuelLog) error {
	filter := bson.M{"_id": fuelLog.ID}
	update := bson.M{
		"$set": bson.M{
			"truck_id":          fuelLog.TruckID,
			"driver_id":         fuelLog.DriverID,
			"date":              fuelLog.Date,
			"gallons_purchased": fuelLog.GallonsPurchased,
			"price_per_gallon":  fuelLog.PricePerGallon,
			"total_cost":        fuelLog.TotalCost,
			"location":          fuelLog.Location,
			"odometer_reading":  fuelLog.OdometerReading,
			"updated_at":        primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.fuelLogs.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update fuel log: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("fuel log entry not found")
	}

	return nil
}

func (r *fuelLogRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.fuelLogs.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete fuel log: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrFuelLogNotFound
	}

	return nil
}

func (r *fuelLogRepository) List(ctx context.Context, limit, offset int64) (*ListFuelLogsResult, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	total, err := r.fuelLogs.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.M{"_id": -1}}},
		{{Key: "$skip", Value: offset}},
		{{Key: "$limit", Value: limit}},
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
		{{
			Key: "$project", Value: bson.M{
				"truck_id":  0,
				"driver_id": 0,
			},
		}},
	}

	cursor, err := r.fuelLogs.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve list of users: %w", err)
	}
	defer cursor.Close(ctx)

	fuelLogs := make([]*domain.FuelLog, 0, limit)
	if err := cursor.All(ctx, &fuelLogs); err != nil {
		return nil, fmt.Errorf("failed to decode fuel logs: %w", err)
	}

	return &ListFuelLogsResult{
		FuelLogs: fuelLogs,
		Total:    total,
	}, nil
}
