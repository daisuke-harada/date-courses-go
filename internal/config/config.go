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
	DB         DBConfig
	GoogleMaps GoogleMapsConfig
	Recruit    RecruitConfig
	Batch      BatchConfig
	JWT        JWTConfig
	CORS       CORSConfig
	RateLimit  RateLimitConfig
}

type GoogleMapsConfig struct {
	APIKey string `envconfig:"GOOGLE_MAPS_API_KEY" required:"true"`
}

type RecruitConfig struct {
	APIKey string `envconfig:"RECRUIT_API_KEY" required:"true"`
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

type RateLimitConfig struct {
	// LoginAttemptsPerMinute はログイン・新規登録の1分あたりの試行上限（IP ごと）です。
	LoginAttemptsPerMinute int `envconfig:"RATE_LIMIT_LOGIN_ATTEMPTS_PER_MINUTE" default:"10"`
}

type CORSConfig struct {
	// AllowOrigins は CORS で許可するオリジンです。カンマ区切りで指定します。
	// 本番のフロントエンドのドメインは環境変数で渡すため、既定値はローカル開発用のみ。
	AllowOrigins []string `envconfig:"CORS_ALLOW_ORIGINS" default:"http://localhost:3000,http://localhost:8080"`
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
		if e := envconfig.Process("", &cfg.DB); e != nil {
			slog.Error("failed to process environment db", "err", e)
		}
		if e := envconfig.Process("", &cfg.GoogleMaps); e != nil {
			slog.Error("failed to process environment google maps", "err", e)
		}
		if e := envconfig.Process("", &cfg.JWT); e != nil {
			slog.Error("failed to process environment jwt", "err", e)
		}
		if e := envconfig.Process("", &cfg.Recruit); e != nil {
			slog.Error("failed to process environment recruit", "err", e)
		}
		if e := envconfig.Process("", &cfg.Batch); e != nil {
			slog.Error("failed to process environment batch", "err", e)
		}
		if e := envconfig.Process("", &cfg.CORS); e != nil {
			slog.Error("failed to process environment cors", "err", e)
		}
		if e := envconfig.Process("", &cfg.RateLimit); e != nil {
			slog.Error("failed to process environment rate limit", "err", e)
		}
	})

	return cfg
}
