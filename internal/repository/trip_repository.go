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
	GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.Trip, error)
	Update(ctx context.Context, trip *domain.Trip) error
	Delete(ctx context.Context, id, userID primitive.ObjectID) error
	List(ctx context.Context, filter domain.TripFilter) (*ListTripsResult, error)
}

type ListTripsResult struct {
	Trips []*domain.Trip
	Total int64
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

func (r *tripRepository) GetById(ctx context.Context, id, userID primitive.ObjectID) (*domain.Trip, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"_id":     id,
			"user_id": userID,
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
	filter := bson.M{
		"_id":     trip.ID,
		"user_id": trip.UserID,
	}

	// First get the existing trip to preserve any fields we don't want to lose
	existingTrip, err := r.GetById(ctx, trip.ID, trip.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing trip: %w", err)
	}

	// Preserve the IDs from the existing trip if they're not being updated
	if trip.DriverID == nil {
		trip.DriverID = existingTrip.DriverID
	}
	if trip.TruckID == nil {
		trip.TruckID = existingTrip.TruckID
	}
	if trip.StartFacilityID == nil {
		trip.StartFacilityID = existingTrip.StartFacilityID
	}
	if trip.EndFacilityID == nil {
		trip.EndFacilityID = existingTrip.EndFacilityID
	}

	update := bson.M{
		"$set": bson.M{
			"trip_number":        trip.TripNumber,
			"driver_id":          trip.DriverID,
			"truck_id":           trip.TruckID,
			"start_facility_id":  trip.StartFacilityID,
			"end_facility_id":    trip.EndFacilityID,
			"departure_time":     trip.DepartureTime,
			"arrival_time":       trip.ArrivalTime,
			"status":             trip.Status,
			"cargo":              trip.Cargo,
			"fuel_usage_gallons": trip.FuelUsage,
			"distance_miles":     trip.DistanceMiles,
			"notes":              trip.Notes,
			"updated_at":         primitive.NewDateTimeFromTime(time.Now()),
			"user_id":            trip.UserID,
		},
	}

	result, err := r.trips.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("trip not found")
	}

	updatedTrip, err := r.GetById(ctx, trip.ID, trip.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch updated trip: %w", err)
	}

	*trip = *updatedTrip

	return nil
}

func (r *tripRepository) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	result, err := r.trips.DeleteOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrTripNotFound
	}

	return nil
}

func (r *tripRepository) List(ctx context.Context, filter domain.TripFilter) (*ListTripsResult, error) {
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

	if filter.DriverID != nil {
		filterQuery["driver_id"] = filter.DriverID
	}

	if filter.TruckID != nil {
		filterQuery["truck_id"] = filter.TruckID
	}

	if filter.StartFacilityID != nil {
		filterQuery["start_facility_id"] = filter.StartFacilityID
	}

	if filter.EndFacilityID != nil {
		filterQuery["end_facility_id"] = filter.EndFacilityID
	}

	total, err := r.trips.CountDocuments(ctx, filterQuery)
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
		return nil, fmt.Errorf("failed to execute aggregate query: %w", err)
	}
	defer cursor.Close(ctx)

	trips := make([]*domain.Trip, 0, filter.Limit)
	if err := cursor.All(ctx, &trips); err != nil {
		return nil, fmt.Errorf("failed to decode trips: %w", err)
	}

	return &ListTripsResult{
		Trips: trips,
		Total: total,
	}, nil
}
