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

type tripRepository struct {
	trips *mongo.Collection
}

type TripRepository interface {
	Create(ctx context.Context, trip *domain.Trip) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.Trip, error)
	Update(ctx context.Context, trip *domain.Trip) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.Trip, error)
}

func NewTripRepository(db *database.MongoDB) TripRepository {
	return &tripRepository{
		trips: db.Database.Collection("trips"),
	}
}

func (r *tripRepository) Create(ctx context.Context, trip *domain.Trip) error {
	now := time.Now()
	trip.CreatedAt = primitive.NewDateTimeFromTime(now)
	trip.UpdatedAt = primitive.NewDateTimeFromTime(now)

	_, err := r.trips.InsertOne(ctx, trip)
	if err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}

	return nil
}

func (r *tripRepository) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Trip, error) {
	filter := bson.M{"_id": id}

	var trip domain.Trip
	err := r.trips.FindOne(ctx, filter).Decode(&trip)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}
	return &trip, nil
}

func (r *tripRepository) Update(ctx context.Context, trip *domain.Trip) error {
	filter := bson.M{"_id": trip.ID}
	update := bson.M{
		"$set": bson.M{
			"trip_number":    trip.TripNumber,
			"driver_id":      trip.DriverID,
			"truck_id":       trip.TruckID,
			"start_facility": trip.StartFacility,
			"end_facility":   trip.EndFacility,
			"route":          trip.Route,
			"start_time":     trip.StartTime,
			"end_time":       trip.EndTime,
			"status":         trip.Status,
			"Cargo":          trip.Cargo,
			"FuelUsage":      trip.FuelUsage,
			"distance_miles": trip.DistanceMiles,
			"updated_at":     primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err := r.trips.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	return nil
}

func (r *tripRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.trips.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}
	return nil
}

func (r *tripRepository) List(ctx context.Context, limit, offset int64) ([]*domain.Trip, error) {
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := r.trips.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve list of users: %w", err)
	}
	defer cursor.Close(ctx)

	var trips []*domain.Trip
	for cursor.Next(ctx) {
		var d domain.Trip
		if err := cursor.Decode(&d); err != nil {
			return nil, fmt.Errorf("failed to decode trip: %w", err)
		}
		trips = append(trips, &d)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return trips, nil
}
