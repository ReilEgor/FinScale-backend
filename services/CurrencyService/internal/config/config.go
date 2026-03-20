package config

type (
	RedisHostType     string
	RedisPortType     string
	RedisPasswordType string

	CompareFinAPIURLType string
	CompareFinAPIKeyType string
)
type Config struct {
	RedisHost        RedisHostType        `env:"REDIS_HOST" envDefault:"redis"`
	RedisPort        RedisPortType        `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword    RedisPasswordType    `env:"REDIS_PASSWORD" envDefault:"redis_password"`
	CompareFinAPIURL CompareFinAPIURLType `env:"COMPARE_FIN_API_URL" envDefault:"https://api.coingecko.com/api/v3"`
	CompareFinAPIKey CompareFinAPIKeyType `env:"COMPARE_FIN_API_KEY" envDefault:"your_api_key"`
}
