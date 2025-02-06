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

type facilityRepository struct {
	facilities *mongo.Collection
}

type FacilityRepository interface {
	Create(ctx context.Context, facility *domain.Facility) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.Facility, error)
	Update(ctx context.Context, facility *domain.Facility) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.Facility, error)
}

func NewFacilityRepository(db *database.MongoDB) FacilityRepository {
	return &facilityRepository{
		facilities: db.Database.Collection("facilities"),
	}
}

func (r *facilityRepository) Create(ctx context.Context, facility *domain.Facility) error {
	now := time.Now()
	facility.CreatedAt = primitive.NewDateTimeFromTime(now)
	facility.UpdatedAt = primitive.NewDateTimeFromTime(now)

	_, err := r.facilities.InsertOne(ctx, facility)
	if err != nil {
		return fmt.Errorf("failed to create facility: %w", err)
	}

	return nil
}

func (r *facilityRepository) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Facility, error) {
	filter := bson.M{"_id": id}

	var facility domain.Facility
	err := r.facilities.FindOne(ctx, filter).Decode(&facility)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get facility: %w", err)
	}
	return &facility, nil
}

func (r *facilityRepository) Update(ctx context.Context, facility *domain.Facility) error {
	filter := bson.M{"_id": facility.ID}
	update := bson.M{
		"$set": bson.M{
			"facility_number":    facility.FacilityNumber,
			"name":               facility.Name,
			"type":               facility.Type,
			"address":            facility.Address,
			"contact_info":       facility.ContactInfo,
			"parking_capacity":   facility.ParkingCapacity,
			"services_available": facility.ServicesAvailable,
			"updated_at":         primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.facilities.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update facility: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("facility not found")
	}

	return nil
}

func (r *facilityRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.facilities.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete facility: %w", err)
	}
	return nil
}

func (r *facilityRepository) List(ctx context.Context, limit, offset int64) ([]*domain.Facility, error) {
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := r.facilities.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve list of users: %w", err)
	}
	defer cursor.Close(ctx)

	var facilities []*domain.Facility
	for cursor.Next(ctx) {
		var d domain.Facility
		if err := cursor.Decode(&d); err != nil {
			return nil, fmt.Errorf("failed to decode facility: %w", err)
		}
		facilities = append(facilities, &d)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return facilities, nil
}
