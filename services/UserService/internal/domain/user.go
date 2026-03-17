package domain

import (
	"context"
	"time"
)

type User struct {
	ID         int       `db:"id" json:"id"`
	KeycloakID string    `db:"keycloak_id" json:"keycloak_id"`
	Username   string    `db:"username" json:"username"`
	Email      string    `db:"email" json:"email"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type UserUseCase interface {
	CreateUser(ctx context.Context, keycloak_id, name, email string) (int64, error)
	GetUser(ctx context.Context, keycloak_id string) (*User, error)
	UpdateUser(ctx context.Context, keycloak_id, name, email string) error
	DeleteUser(ctx context.Context, id int64) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, keycloak_id, name, email string) (int64, error)
	GetUser(ctx context.Context, keycloak_id string) (*User, error)
	UpdateUser(ctx context.Context, keycloak_id, name, email string) error
	DeleteUser(ctx context.Context, id int64) error
}
