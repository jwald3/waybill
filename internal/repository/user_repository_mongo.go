package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jwald3/go_rest_template/internal/database"
	"github.com/jwald3/go_rest_template/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userMongoRepository struct {
	collection *mongo.Collection
}

type UserMongoRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

func NewUserMongoRepository(mdb *database.MongoDB) UserMongoRepository {
	return &userMongoRepository{
		collection: mdb.Database.Collection("users"),
	}
}

func (r *userMongoRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user in mongo: %w", err)
	}

	return nil
}

func (r *userMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	filter := bson.M{"_id": id}

	var user domain.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user from mongo: %w", err)
	}
	return &user, nil
}

func (r *userMongoRepository) Update(ctx context.Context, user *domain.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"email":         user.Email,
			"password_hash": user.Password.Hash(),
			"status":        user.Status,
			"updated_at":    time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user in mongo: %w", err)
	}

	return nil
}

func (r *userMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user in mongo: %w", err)
	}
	return nil
}

func (r *userMongoRepository) List(ctx context.Context, limit, offset int64) ([]*domain.User, error) {
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list users in mongo: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var u domain.User
		if err := cursor.Decode(&u); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, &u)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}

func (r *userMongoRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to count documents: %w", err)
	}
	return count > 0, nil
}
