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

type truckRepository struct {
	trucks *mongo.Collection
}

type TruckRepository interface {
	Create(ctx context.Context, truck *domain.Truck) error
	GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.Truck, error)
	Update(ctx context.Context, truck *domain.Truck) error
	Delete(ctx context.Context, id, userID primitive.ObjectID) error
	List(ctx context.Context, filter domain.TruckFilter) (*ListTrucksResult, error)
}

type ListTrucksResult struct {
	Trucks []*domain.Truck
	Total  int64
}

func NewTruckRepository(db *database.MongoDB) TruckRepository {
	return &truckRepository{
		trucks: db.Database.Collection("trucks"),
	}
}

func (r *truckRepository) Create(ctx context.Context, truck *domain.Truck) error {
	now := time.Now()
	truck.CreatedAt = primitive.NewDateTimeFromTime(now)
	truck.UpdatedAt = primitive.NewDateTimeFromTime(now)

	_, err := r.trucks.InsertOne(ctx, truck)
	if err != nil {
		return fmt.Errorf("failed to create truck: %w", err)
	}

	return nil
}

func (r *truckRepository) GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.Truck, error) {
	pipeline := mongo.Pipeline{
		{{
			Key: "$match", Value: bson.M{
				"_id":     id,
				"user_id": userID,
			},
		}},
		{{
			Key: "$lookup", Value: bson.M{
				"from":         "drivers",
				"localField":   "assigned_driver_id",
				"foreignField": "_id",
				"as":           "assigned_driver",
			},
		}},
		{{
			Key: "$unwind", Value: bson.M{
				"path":                       "$assigned_driver",
				"preserveNullAndEmptyArrays": true,
			},
		}},
		{{
			Key: "$project", Value: bson.M{
				"assigned_driver_id": 0,
			},
		}},
	}

	var result domain.Truck
	cursor, err := r.trucks.Aggregate(ctx, pipeline)
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
		return nil, fmt.Errorf("failed to decode truck: %w", err)
	}

	return &result, nil
}

func (r *truckRepository) Update(ctx context.Context, truck *domain.Truck) error {
	filter := bson.M{"_id": truck.ID}
	update := bson.M{
		"$set": bson.M{
			"mileage":          truck.Mileage,
			"status":           truck.Status,
			"last_maintenance": truck.LastMaintenance,
			"updated_at":       primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.trucks.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update truck: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("truck not found")
	}

	return nil
}

func (r *truckRepository) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	result, err := r.trucks.DeleteOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete truck: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrTruckNotFound
	}

	return nil
}

func (r *truckRepository) List(ctx context.Context, filter domain.TruckFilter) (*ListTrucksResult, error) {
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

	if filter.UserID != primitive.NilObjectID {
		filterQuery["user_id"] = filter.UserID
	}

	if filter.TrailerType != "" {
		filterQuery["trailer_type"] = filter.TrailerType
	}

	if filter.FuelType != "" {
		filterQuery["fuel_type"] = filter.FuelType
	}

	if filter.Status != "" {
		filterQuery["status"] = filter.Status
	}

	if filter.AssignedDriverID != &primitive.NilObjectID {
		filterQuery["assigned_driver_id"] = filter.AssignedDriverID
	}

	total, err := r.trucks.CountDocuments(ctx, filterQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filterQuery}},
		{{Key: "$sort", Value: bson.M{"_id": -1}}},
		{{Key: "$skip", Value: filter.Offset}},
		{{Key: "$limit", Value: filter.Limit}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "drivers",
			"localField":   "assigned_driver_id",
			"foreignField": "_id",
			"as":           "assigned_driver",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$assigned_driver",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"assigned_driver_id": 0,
		}}},
	}

	cursor, err := r.trucks.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate query: %w", err)
	}
	defer cursor.Close(ctx)

	trucks := make([]*domain.Truck, 0, filter.Limit)
	if err := cursor.All(ctx, &trucks); err != nil {
		return nil, fmt.Errorf("failed to decode trucks: %w", err)
	}

	return &ListTrucksResult{
		Trucks: trucks,
		Total:  total,
	}, nil
}
