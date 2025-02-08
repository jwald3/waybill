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
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": id}}},
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
			"from":         "facilities",
			"localField":   "start_facility_id",
			"foreignField": "_id",
			"as":           "start_facility",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$start_facility",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "facilities",
			"localField":   "end_facility_id",
			"foreignField": "_id",
			"as":           "end_facility",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$end_facility",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"driver_id":         0,
			"truck_id":          0,
			"start_facility_id": 0,
			"end_facility_id":   0,
		}}},
	}

	var result domain.Trip
	cursor, err := r.trips.Aggregate(ctx, pipeline)
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
		return nil, fmt.Errorf("failed to decode trip: %w", err)
	}

	return &result, nil
}

func (r *tripRepository) Update(ctx context.Context, trip *domain.Trip) error {
	filter := bson.M{"_id": trip.ID}
	update := bson.M{
		"$set": bson.M{
			"trip_number":       trip.TripNumber,
			"driver_id":         trip.DriverID,
			"truck_id":          trip.TruckID,
			"start_facility_id": trip.StartFacilityID,
			"end_facility_id":   trip.EndFacilityID,
			"route":             trip.Route,
			"start_time":        trip.StartTime,
			"end_time":          trip.EndTime,
			"status":            trip.Status,
			"Cargo":             trip.Cargo,
			"FuelUsage":         trip.FuelUsage,
			"distance_miles":    trip.DistanceMiles,
			"updated_at":        primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.trips.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("trip not found")
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
	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.M{"_id": -1}}},
		{{Key: "$skip", Value: offset}},
		{{Key: "$limit", Value: limit}},
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
			"from":         "facilities",
			"localField":   "start_facility_id",
			"foreignField": "_id",
			"as":           "start_facility",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$start_facility",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "facilities",
			"localField":   "end_facility_id",
			"foreignField": "_id",
			"as":           "end_facility",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$end_facility",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"driver_id":         0,
			"truck_id":          0,
			"start_facility_id": 0,
			"end_facility_id":   0,
		}}},
	}

	cursor, err := r.trips.Aggregate(ctx, pipeline)
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
