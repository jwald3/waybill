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
	ListWithFilter(ctx context.Context, filter domain.FacilityFilter) (*ListFacilitiesResult, error)
	UpdateAvailableFacilityServices(ctx context.Context, id primitive.ObjectID, servicesAvailable []domain.FacilityService) error
}

type ListFacilitiesResult struct {
	Facilities []*domain.Facility
	Total      int64
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
	result, err := r.facilities.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete facility: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrFacilityNotFound
	}

	return nil
}

func (r *facilityRepository) ListWithFilter(ctx context.Context, filter domain.FacilityFilter) (*ListFacilitiesResult, error) {
	// handle wrangling the filter values to make sure they're within the bounds we want to support.
	// if we don't do this, we may get some unexpected results (including requesting 1000+ results from the db)
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

	// go through each of the filter options and add them to the filter query if they're not empty
	if filter.StateCode != "" {
		filterQuery["address.state"] = filter.StateCode
	}

	if filter.Type != "" {
		filterQuery["type"] = filter.Type
	}

	if len(filter.ServicesInclude) > 0 {
		filterQuery["services_available"] = bson.M{
			"$all": filter.ServicesInclude,
		}
	}

	if filter.MinCapacity != nil {
		filterQuery["parking_capacity"] = bson.M{
			"$gte": *filter.MinCapacity,
		}
	}

	if filter.MaxCapacity != nil {
		if _, exists := filterQuery["parking_capacity"]; exists {
			filterQuery["parking_capacity"].(bson.M)["$lte"] = *filter.MaxCapacity
		} else {
			filterQuery["parking_capacity"] = bson.M{
				"$lte": *filter.MaxCapacity,
			}
		}
	}

	// find the total number of facilities that match the filter (this is prior to pagination but after the filter is applied,
	// so we're counting only the facilities that match the filter)
	total, err := r.facilities.CountDocuments(ctx, filterQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get total count: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(filter.Limit)
	findOptions.SetSkip(filter.Offset)
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	// find the facilities that match the filter and return paginated results
	cursor, err := r.facilities.Find(ctx, filterQuery, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve list of facilities: %w", err)
	}
	defer cursor.Close(ctx)

	facilities := make([]*domain.Facility, 0, filter.Limit)
	if err := cursor.All(ctx, &facilities); err != nil {
		return nil, fmt.Errorf("failed to decode facilities: %w", err)
	}

	return &ListFacilitiesResult{
		Facilities: facilities,
		Total:      total,
	}, nil
}

func (r *facilityRepository) UpdateAvailableFacilityServices(ctx context.Context, id primitive.ObjectID, servicesAvailable []domain.FacilityService) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"services_available": servicesAvailable,
			"updated_at":         primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.facilities.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update facility services: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrFacilityNotFound
	}

	return nil
}
