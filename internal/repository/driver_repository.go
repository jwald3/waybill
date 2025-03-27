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

type driverRepository struct {
	drivers *mongo.Collection
}

type DriverRepository interface {
	Create(ctx context.Context, driver *domain.Driver) error
	GetById(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*domain.Driver, error)
	Update(ctx context.Context, driver *domain.Driver) error
	Delete(ctx context.Context, id, userID primitive.ObjectID) error
	List(ctx context.Context, filter domain.DriverFilter) (*ListDriversResult, error)
	UpdateEmploymentStatus(ctx context.Context, id primitive.ObjectID, status domain.EmploymentStatus) error
}

type ListDriversResult struct {
	Drivers []*domain.Driver
	Total   int64
}

type DriverFilter struct {
	UserID           primitive.ObjectID
	LicenseState     string
	Phone            domain.PhoneNumber
	Email            domain.Email
	EmploymentStatus domain.EmploymentStatus
	Limit            int64
	Offset           int64
}

func NewDriverRepository(db *database.MongoDB) DriverRepository {
	return &driverRepository{
		drivers: db.Database.Collection("drivers"),
	}
}

func (r *driverRepository) Create(ctx context.Context, driver *domain.Driver) error {
	now := time.Now()
	driver.CreatedAt = primitive.NewDateTimeFromTime(now)
	driver.UpdatedAt = primitive.NewDateTimeFromTime(now)

	_, err := r.drivers.InsertOne(ctx, driver)
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	return nil
}

func (r *driverRepository) GetById(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*domain.Driver, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match",
			Value: bson.M{
				"_id":     id,
				"user_id": userID,
			},
		}},
		{{
			Key: "$lookup", Value: bson.M{
				"from":         "trucks",
				"localField":   "assigned_truck_id",
				"foreignField": "_id",
				"as":           "assigned_truck",
			},
		}},
		{{
			Key: "$unwind", Value: bson.M{
				"path":                       "$assigned_truck",
				"preserveNullAndEmptyArrays": true,
			},
		}},
		{{
			Key: "$project", Value: bson.M{
				"assigned_truck_id": 0,
			},
		}},
	}

	var result domain.Driver
	cursor, err := r.drivers.Aggregate(ctx, pipeline)
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
		return nil, fmt.Errorf("failed to decode driver: %w", err)
	}

	return &result, err
}

func (r *driverRepository) Update(ctx context.Context, driver *domain.Driver) error {
	filter := bson.M{"_id": driver.ID}
	update := bson.M{
		"$set": bson.M{
			"first_name":        driver.FirstName,
			"last_name":         driver.LastName,
			"phone":             driver.Phone,
			"email":             driver.Email,
			"address":           driver.Address,
			"employment_status": driver.EmploymentStatus,
			"updated_at":        primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.drivers.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("driver not found")
	}

	return nil
}

func (r *driverRepository) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	result, err := r.drivers.DeleteOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrDriverNotFound
	}

	return nil
}

func (r *driverRepository) List(ctx context.Context, filter domain.DriverFilter) (*ListDriversResult, error) {
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

	if filter.LicenseState != "" {
		filterQuery["license_state"] = filter.LicenseState
	}
	if filter.Phone != "" {
		filterQuery["phone"] = filter.Phone
	}
	if filter.Email != "" {
		filterQuery["email"] = filter.Email
	}
	if filter.EmploymentStatus != "" {
		filterQuery["employment_status"] = filter.EmploymentStatus
	}

	total, err := r.drivers.CountDocuments(ctx, filterQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filterQuery}},
		{{Key: "$sort", Value: bson.M{"_id": -1}}},
		{{Key: "$skip", Value: filter.Offset}},
		{{Key: "$limit", Value: filter.Limit}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "trucks",
			"localField":   "assigned_truck_id",
			"foreignField": "_id",
			"as":           "assigned_truck",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$assigned_truck",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{
			Key: "$project", Value: bson.M{
				"assigned_truck_id": 0,
			},
		}},
	}

	cursor, err := r.drivers.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate query: %w", err)
	}
	defer cursor.Close(ctx)

	drivers := make([]*domain.Driver, 0, filter.Limit)
	if err := cursor.All(ctx, &drivers); err != nil {
		return nil, fmt.Errorf("failed to decode drivers: %w", err)
	}

	return &ListDriversResult{
		Drivers: drivers,
		Total:   total,
	}, nil
}

func (r *driverRepository) UpdateEmploymentStatus(ctx context.Context, id primitive.ObjectID, status domain.EmploymentStatus) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"employment_status": status,
			"updated_at":        primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.drivers.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update driver employment status: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrDriverNotFound
	}

	return nil
}
