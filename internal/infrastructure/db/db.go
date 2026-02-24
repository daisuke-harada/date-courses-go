package db

import (
	"context"
	"fmt"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connector interface {
	Open(ctx context.Context, cfg config.DBConfig) (*gorm.DB, error)
}

type PostgresConnector struct{}

func (PostgresConnector) Open(ctx context.Context, cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	gdb, err := gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return gdb, nil
}

func NewConnector() Connector {
	return PostgresConnector{}
}

// Connect uses a Connector to open a *gorm.DB. Currently delegates to PostgresConnector.
func Connect(ctx context.Context, cfg config.DBConfig) (*gorm.DB, error) {
	connector := NewConnector()
	return connector.Open(ctx, cfg)
}
