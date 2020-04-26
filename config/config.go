package config

import (
	"github.com/caarlos0/env/v6"
)

//"port=5432 host=localhost user=postgres password=111111 dbname=postgres sslmode=disable
type Config struct {
	Listen string `env:"LISTEN" envDefault:":8080"`
	PgSQL  string `env:"PGSQL" envDefault:"postgres://postgres:111111@localhost:5432/postgres?sslmode=disable"`
	//port=5432 host=localhost user=postgres password=111111 dbname=postgres sslmode=disable
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
}

// Load config
func Load() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
