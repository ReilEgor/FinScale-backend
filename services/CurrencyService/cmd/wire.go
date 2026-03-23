//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/api"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/repository/redis"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/transport/rest"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/transport/rest/handlers"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/usecase"
	"github.com/google/wire"
)

var APISet = wire.NewSet(
	api.NewCryptoCompare,
	wire.Bind(new(domain.CurrencyFetcher), new(*api.CryptoCompare)),
)

var UsecaseSet = wire.NewSet(
	usecase.NewCurrencyUseCase,
	wire.Bind(new(domain.CurrencyUseCase), new(*usecase.CurrencyUseCase)),
)

var RestSet = wire.NewSet(
	rest.NewGinServer,
	handlers.NewHandler,
)

var RedisSet = wire.NewSet(
	redis.NewRedisClient,
	redis.NewCurrencyRepository,
	wire.Bind(new(domain.CurrencyRepository), new(*redis.CurrencyRepository)),
)

type App struct {
	Logic  domain.CurrencyUseCase
	Server *rest.GinServer
}

func InitializeApp(
	redisHost config.RedisHostType,
	redisPort config.RedisPortType,
	redisPassword config.RedisPasswordType,
	redisDB int,
	coinGeckoAPIURL config.CompareFinAPIURLType,
	coinGeckoAPIKeyType config.CompareFinAPIKeyType,
) (*App, func(), error) {
	wire.Build(
		APISet,
		RedisSet,
		UsecaseSet,
		RestSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
