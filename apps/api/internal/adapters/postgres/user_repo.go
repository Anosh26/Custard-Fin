package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) ports.UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (id, username, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return fmt.Errorf("%w: user with username %s already exists", domain.ErrConflict, user.Username)
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	username = strings.TrimSpace(username)
	query := `
	SELECT id, username, password_hash, created_at, updated_at
	FROM users
	WHERE username = $1
	`
	row := r.pool.QueryRow(ctx, query, username)

	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to query user by username: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	id = strings.TrimSpace(id)
	query := `
	SELECT id, username, password_hash, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to query user by id: %w", err)
	}
	return &u, nil
}
