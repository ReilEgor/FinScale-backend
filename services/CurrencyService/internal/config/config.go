package config

type (
	RedisHostType     string
	RedisPortType     string
	RedisPasswordType string
)
type Config struct {
	RedisHost     RedisHostType     `env:"REDIS_HOST" envDefault:"redis"`
	RedisPort     RedisPortType     `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword RedisPasswordType `env:"REDIS_PASSWORD" envDefault:"redis_password"`
}
