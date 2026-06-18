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
	DB           DBConfig
	GoogleMaps   GoogleMapsConfig
	GooglePlaces GooglePlacesConfig
	Gemini       GeminiConfig
	Batch        BatchConfig
	JWT          JWTConfig
}

type GoogleMapsConfig struct {
	APIKey string `envconfig:"GOOGLE_MAPS_API_KEY" required:"true"`
}

type GooglePlacesConfig struct {
	APIKey string `envconfig:"GOOGLE_PLACES_API_KEY" required:"true"`
}

type GeminiConfig struct {
	APIKey string `envconfig:"GEMINI_API_KEY" required:"true"`
	Model  string `envconfig:"GEMINI_MODEL" default:"gemini-2.0-flash-lite"`
}

type BatchConfig struct {
	SpotsPerCombination  int `envconfig:"BATCH_SPOTS_PER_COMBINATION" default:"5"`
	MinExistingSpots     int `envconfig:"BATCH_MIN_EXISTING_SPOTS" default:"5"`
	MaxTasksPerRun       int `envconfig:"BATCH_MAX_TASKS_PER_RUN" default:"50"`
	MaxRequestsPerMinute int `envconfig:"BATCH_MAX_REQUESTS_PER_MINUTE" default:"60"`
}

type JWTConfig struct {
	SecretKey string `envconfig:"JWT_SECRET_KEY" required:"true"`
}

type DBConfig struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     int    `envconfig:"DB_PORT" default:"3306"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
	TLS      bool   `envconfig:"DB_TLS" default:"false"`
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
		if e := envconfig.Process("", &cfg.GoogleMaps); e != nil {
			slog.Error("failed to process environment google maps", "err", e)
		}
		if e := envconfig.Process("", &cfg.JWT); e != nil {
			slog.Error("failed to process environment jwt", "err", e)
		}
		if e := envconfig.Process("", &cfg.GooglePlaces); e != nil {
			slog.Error("failed to process environment google places", "err", e)
		}
		if e := envconfig.Process("", &cfg.Gemini); e != nil {
			slog.Error("failed to process environment gemini", "err", e)
		}
		if e := envconfig.Process("", &cfg.Batch); e != nil {
			slog.Error("failed to process environment batch", "err", e)
		}
	})

	return cfg
}
