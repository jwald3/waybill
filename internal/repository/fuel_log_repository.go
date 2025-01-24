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

type fuelLogRepository struct {
	fuelLogs *mongo.Collection
}

type FuelLogRepository interface {
	Create(ctx context.Context, fuelLog *domain.FuelLog) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.FuelLog, error)
	Update(ctx context.Context, fuelLog *domain.FuelLog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.FuelLog, error)
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
	filter := bson.M{"_id": id}

	var fuelLog domain.FuelLog
	err := r.fuelLogs.FindOne(ctx, filter).Decode(&fuelLog)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fuelLog: %w", err)
	}
	return &fuelLog, nil
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

	_, err := r.fuelLogs.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update fuelLog: %w", err)
	}

	return nil
}

func (r *fuelLogRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.fuelLogs.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete fuelLog: %w", err)
	}
	return nil
}

func (r *fuelLogRepository) List(ctx context.Context, limit, offset int64) ([]*domain.FuelLog, error) {
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := r.fuelLogs.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve list of users: %w", err)
	}
	defer cursor.Close(ctx)

	var fuelLogs []*domain.FuelLog
	for cursor.Next(ctx) {
		var d domain.FuelLog
		if err := cursor.Decode(&d); err != nil {
			return nil, fmt.Errorf("failed to decode fuelLog: %w", err)
		}
		fuelLogs = append(fuelLogs, &d)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return fuelLogs, nil
}
