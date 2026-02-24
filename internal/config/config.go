package config

import (
	"log/slog"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	once sync.Once
	cfg  *Config
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
	// Connection pool settings
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"25"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"5m"`
}

func Get() *Config {
	once.Do(func() {
		cfg = &Config{}
		// Use empty prefix because DB field tags already contain full env names (e.g. DB_HOST).
		if e := envconfig.Process("", &cfg.DB); e != nil {
			// use conventional key name "err"
			slog.Error("failed to process environment db", "err", e)
		}
	})

	return cfg
}
