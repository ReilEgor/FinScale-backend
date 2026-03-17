package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/UserService/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: slog.With(slog.String("component", "UserRepository")),
	}
}

const createUserQuery = `
	INSERT INTO users (keycloak_id, username, email)
	VALUES ($1, $2, $3) RETURNING id;
`

func (r *UserRepository) CreateUser(ctx context.Context, keycloak_id, name, email string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, createUserQuery, keycloak_id, name, email).Scan(&id)
	if err != nil {
		r.logger.Error("failed to create user", slog.Any("error", err))
		return id, err
	}
	return id, nil
}

const getUserQuery = `
	SELECT id, keycloak_id, username, email, created_at FROM users WHERE keycloak_id = $1
`

func (r *UserRepository) GetUser(ctx context.Context, keycloak_id string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(ctx, getUserQuery, keycloak_id).Scan(
		&user.ID,
		&user.KeycloakID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("failed to get user", slog.Any("error", err))
		return nil, err
	}
	return user, nil
}

const updateUserQuery = `
    UPDATE users 
    SET username = $2, 
        email = $3 
    WHERE keycloak_id = $1
    RETURNING id;
`

func (r *UserRepository) UpdateUser(ctx context.Context, keycloak_id, name, email string) error {
	result, err := r.db.Exec(ctx, updateUserQuery, keycloak_id, name, email)
	if err != nil {
		r.logger.Error("pgpool: update failed", slog.String("id", keycloak_id), slog.Any("err", err))
		return err
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user with keycloak_id %s not found", keycloak_id)
	}

	return nil
}

const deleteUserQuery = `
    DELETE FROM users 
    WHERE id = $1
`

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	result, err := r.db.Exec(ctx, deleteUserQuery, id)
	if err != nil {
		r.logger.Error("failed to delete user", slog.Int64("id", id), slog.Any("error", err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Warn("no user found to delete", slog.Int64("id", id))
	}

	return nil
}
