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

type driverRepository struct {
	drivers *mongo.Collection
}

type DriverRepository interface {
	Create(ctx context.Context, driver *domain.Driver) error
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.Driver, error)
	Update(ctx context.Context, driver *domain.Driver) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.Driver, error)
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

func (r *driverRepository) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Driver, error) {
	filter := bson.M{"_id": id}

	var driver domain.Driver
	err := r.drivers.FindOne(ctx, filter).Decode(&driver)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}
	return &driver, nil
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

	_, err := r.drivers.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}

	return nil
}

func (r *driverRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.drivers.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}
	return nil
}

func (r *driverRepository) List(ctx context.Context, limit, offset int64) ([]*domain.Driver, error) {
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := r.drivers.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve list of users: %w", err)
	}
	defer cursor.Close(ctx)

	var drivers []*domain.Driver
	for cursor.Next(ctx) {
		var d domain.Driver
		if err := cursor.Decode(&d); err != nil {
			return nil, fmt.Errorf("failed to decode driver: %w", err)
		}
		drivers = append(drivers, &d)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return drivers, nil
}
