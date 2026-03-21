//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/config"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/repository/postgres"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/transport/rest"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/transport/rest/handlers"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/usecase"
	"github.com/google/wire"
)

var UseCaseSet = wire.NewSet(
	usecase.NewTransactionUseCase,
	wire.Bind(new(domain.TransactionUseCase), new(*usecase.TransactionUseCase)),
)

var RepositorySet = wire.NewSet(
	postgres.NewPostgresRepository,
	postgres.NewTransactionRepository,
	wire.Bind(new(domain.TransactionRepository), new(*postgres.TransactionRepository)),
)

var RestSet = wire.NewSet(
	rest.NewGinServer,
	handlers.NewHandler,
)

type App struct {
	Logic  domain.TransactionUseCase
	Server *rest.GinServer
}

func InitializeApp(ctx context.Context, dsn config.DSN) (*App, func(), error) {
	wire.Build(
		RepositorySet,
		UseCaseSet,
		RestSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
