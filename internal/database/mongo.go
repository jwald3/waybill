package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jwald3/go_rest_template/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoConnection(cfg config.Config) (*MongoDB, error) {
	mongoDSN := cfg.GetDSN()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoDSN)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	dbName := cfg.Database.DBName

	mongoDB := &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}

	if err := mongoDB.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	return mongoDB, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := m.Client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongo health check failed: %w", err)
	}

	return nil
}

func (m *MongoDB) ExecuteTx(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	// start a session
	session, err := m.Client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (any, error) {
		return nil, fn(sessCtx)
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}
