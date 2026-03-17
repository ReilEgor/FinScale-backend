//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/ReilEgor/FinScale-backend/UserService/internal/domain"
	postgres2 "github.com/ReilEgor/FinScale-backend/UserService/internal/repository/postgres"
	"github.com/ReilEgor/FinScale-backend/UserService/internal/transport/rest"
	"github.com/ReilEgor/FinScale-backend/UserService/internal/transport/rest/handlers"
	"github.com/ReilEgor/FinScale-backend/UserService/internal/usecase"
	"github.com/google/wire"
)

var UseCaseSet = wire.NewSet(
	usecase.NewUserUseCase,
	wire.Bind(new(domain.UserUseCase), new(*usecase.UserUseCase)),
)

var RepositorySet = wire.NewSet(
	postgres2.NewPostgresRepository,
	postgres2.NewUserRepository,
	wire.Bind(new(domain.UserRepository), new(*postgres2.UserRepository)),
)

var RestSet = wire.NewSet(
	rest.NewGinServer,
	handlers.NewHandler,
)

type App struct {
	Logic  domain.UserUseCase
	Server *rest.GinServer
}

func InitializeApp(ctx context.Context, dsn string) (*App, func(), error) {
	wire.Build(
		RepositorySet,
		UseCaseSet,
		RestSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
