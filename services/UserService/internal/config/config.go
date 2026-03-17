package config

type (
	DSN string
)

type Config struct {
	DSN DSN `env:"DB_SOURCE" envDefault:"postgres://user:password@localhost:5432/userservice?sslmode=disable"`
}
