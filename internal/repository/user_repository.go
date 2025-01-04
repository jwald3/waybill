package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jwald3/go_rest_template/internal/database"
	"github.com/jwald3/go_rest_template/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

// we're using transaction and non-transaction versions of methods that could entail using transactions.
// these are typically going to be the operations that write to the database.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	CreateTx(ctx context.Context, tx *database.Transaction, user *domain.User) error

	GetByID(ctx context.Context, id int64) (*domain.User, error)

	Update(ctx context.Context, user *domain.User) error
	UpdateTx(ctx context.Context, tx *database.Transaction, user *domain.User) error

	Delete(ctx context.Context, id int64) error
	DeleteTx(ctx context.Context, tx *database.Transaction, id int64) error

	List(ctx context.Context, limit, offset int) ([]*domain.User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByEmailTx(ctx context.Context, tx *database.Transaction, email string) (bool, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	// just basic raw SQL with $ params. ensure that the syntax matches your database engine (this template assunes you're using postgres)
	query := `
		INSERT INTO users (email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()

	// this is where we differ between the transaction and non-transaction functions. functions that do not require inclusion in a transaction
	// can call the database directly. In either case, we populate those parameters with QueryRowContext
	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.Password.Hash(),
		user.Status,
		now,
		now,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var (
		user         domain.User
		passwordHash string
	)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&passwordHash,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Password = domain.NewPasswordFromHash(passwordHash)
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Password.Hash(),
		user.Status,
		now,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.UpdatedAt = now
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, status, created_at, updated_at
		FROM users
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var (
			u            domain.User
			passwordHash string
		)
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&passwordHash,
			&u.Status,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		u.Password = domain.NewPasswordFromHash(passwordHash)
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

// helper method to ensure that you can make basic business logic checks without needing to rely on database constraints
// you can add in more checks like this to include in the service methods
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE email = $1
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// transaction methods:
// these methods are the exact same as their non-transaction counterparts but use the transasction instead of a
// direct reference to the database when making the SQL call. we can optionally use an interface to eliminate the
// need for separate methods, but this redundancy may be easier to grasp when starting.
func (r *userRepository) CreateTx(ctx context.Context, tx *database.Transaction, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	err := tx.QueryRowContext(ctx, query,
		user.Email,
		user.Password.Hash(),
		user.Status,
		now,
		now,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *userRepository) UpdateTx(ctx context.Context, tx *database.Transaction, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	now := time.Now()
	_, err := tx.ExecContext(ctx, query,
		user.Email,
		user.Password.Hash(),
		user.Status,
		now,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.UpdatedAt = now
	return nil
}

func (r *userRepository) DeleteTx(ctx context.Context, tx *database.Transaction, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepository) ExistsByEmailTx(ctx context.Context, tx *database.Transaction, email string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE email = $1
	`
	var count int
	err := tx.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}
