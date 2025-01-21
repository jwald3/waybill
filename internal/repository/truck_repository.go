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

type truckRepository struct {
	trucks *mongo.Collection
}

type TruckRepository interface {
	Create(ctx context.Context, truck *domain.Truck) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.Truck, error)
	Update(ctx context.Context, truck *domain.Truck) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.Truck, error)
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

func (r *truckRepository) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Truck, error) {
	filter := bson.M{"_id": id}

	var truck domain.Truck
	err := r.trucks.FindOne(ctx, filter).Decode(&truck)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get truck: %w", err)
	}

	return &truck, nil
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

	_, err := r.trucks.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update truck: %w", err)
	}

	return nil
}

func (r *truckRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.trucks.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete truck: %w", err)
	}
	return nil
}

func (r *truckRepository) List(ctx context.Context, limit, offset int64) ([]*domain.Truck, error) {
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := r.trucks.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve list of trucks: %w", err)
	}
	defer cursor.Close(ctx)

	var trucks []*domain.Truck
	for cursor.Next(ctx) {
		var t domain.Truck
		if err := cursor.Decode(&t); err != nil {
			return nil, fmt.Errorf("failed to decode truck: %w", err)
		}
		trucks = append(trucks, &t)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return trucks, nil
}
