package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jwald3/go_rest_template/internal/config"
	"github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// utilize the config file to load an instance of DB with defined settings
func NewPostgresConnection(cfg config.Config) (*DB, error) {
	dsn := cfg.GetDSN()

	// uses the DSN we build using the config class to open the connection with the postgres engine. the first arg will be the engine you're using (postgres, mysql, etc.)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// use all of the predefined values as found in the config file
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// this bit is a "sanity check" to make sure that the connection works.
	// If the configurations aren't correct, this will throw an error.
	// Feel free to remove or comment out after ensuring the connection works
	database := &DB{db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	// end sanity check

	return &DB{db}, nil
}

type Transaction struct {
	*sql.Tx
}

// creates a transaction used
func (db *DB) BeginTx(ctx context.Context) (*Transaction, error) {
	tx, err := db.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	return &Transaction{tx}, nil
}

func (db *DB) ExecuteTx(ctx context.Context, fn func(*Transaction) error) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing: %w", err)
	}

	return nil
}

func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (db *DB) PrepareNamedQuery(query string, params map[string]any) (string, []any, error) {
	paramCount := 1
	var values []any

	for key, value := range params {
		param := fmt.Sprintf(":%s", key)
		query = strings.Replace(query, param, fmt.Sprintf("$%d", paramCount), 1)
		values = append(values, value)
		paramCount++
	}

	return query, values, nil
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code {
			case "23505": // unique constraint
				return nil, errors.New("duplicate key value violates unique constraint")
			case "23503": // foreign key contraint
				return nil, errors.New("foreign key constraint")
			case "23502": // not null exception
				return nil, errors.New("null value violated not-null constraint")
			}
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return result, nil
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			return nil, formatPostgresError(pgErr)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return rows, nil
}

// leverages `DB.QueryRowContext` to populate a SQL query string by forwarding a variadic number of arguments.
// we
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}

// a (non-exhaustive) mapper between common postgres error codes and a neatly-formatted error message
// you can find more error codes here: https://www.postgresql.org/docs/current/errcodes-appendix.html
// (skipped codes are uncommon errors or errors that may not be applicable based on operations used)
func formatPostgresError(err *pq.Error) error {
	switch err.Code {
	case "23505":
		return fmt.Errorf("duplicate entry: %s", err.Detail)
	case "23503":
		return fmt.Errorf("foreign key violation: %s", err.Detail)
	case "23502":
		return fmt.Errorf("missing required field: %s", err.Detail)
	default:
		return fmt.Errorf("database error: %s", err.Message)
	}
}
