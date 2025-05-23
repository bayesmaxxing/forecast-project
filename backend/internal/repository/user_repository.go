package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ValidateUser(ctx context.Context, id int64) (bool, error)
	UpdatePassword(ctx context.Context, id int64, password string) error
	ListUsers(ctx context.Context) ([]*models.User, error)
	CreateUserWithPassword(ctx context.Context, username string, passwordHash string) (int64, error)
}

// PostgresUserRepository implements the UserRepository interface
type PostgresUserRepository struct {
	db *database.DB
}

// NewUserRepository creates a new PostgresUserRepository instance
func NewUserRepository(db *database.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()

	query := `INSERT INTO users (username, password, created)
              VALUES ($1, $2, $3)	
              RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Password,
		user.CreatedAt).Scan(&user.ID)
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, username, password, created
              FROM users
              WHERE id = $1`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password, created
              FROM users
              WHERE username = $1`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *PostgresUserRepository) ValidateUser(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := `UPDATE users SET password = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, password)
	return err
}

func (r *PostgresUserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, username, created FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, rows.Err()
}

func (r *PostgresUserRepository) CreateUserWithPassword(ctx context.Context, username string, passwordHash string) (int64, error) {
	createdAt := time.Now()
	var userID int64

	// Direct SQL with no struct involved
	query := `INSERT INTO users (username, password, created)
              VALUES ($1, $2, $3)  
              RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		username,
		passwordHash,
		createdAt).Scan(&userID)

	return userID, err
}
