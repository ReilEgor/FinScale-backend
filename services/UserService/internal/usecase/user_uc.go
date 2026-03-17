package usecase

import (
	"context"
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/UserService/internal/domain"
)

type UserUseCase struct {
	logger *slog.Logger
	db     domain.UserRepository
}

func NewUserUseCase(db domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		db:     db,
		logger: slog.With(slog.String("useCase", "NewUserUseCase")),
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, keycloak_id, name, email string) (int64, error) {
	return uc.db.CreateUser(ctx, keycloak_id, name, email)
}

func (uc *UserUseCase) GetUser(ctx context.Context, keycloak_id string) (*domain.User, error) {
	return uc.db.GetUser(ctx, keycloak_id)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, keycloak_id, name, email string) error {
	return uc.db.UpdateUser(ctx, keycloak_id, name, email)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id int64) error {
	return uc.db.DeleteUser(ctx, id)
}
