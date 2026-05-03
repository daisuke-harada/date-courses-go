package db

import (
	"context"
	"fmt"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connector interface {
	Open(ctx context.Context, cfg config.DBConfig) (*gorm.DB, error)
}

type MySQLConnector struct{}

func (MySQLConnector) Open(ctx context.Context, cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
	return MySQLConnector{}
}

func Connect(ctx context.Context, cfg config.DBConfig) (*gorm.DB, error) {
	connector := NewConnector()
	return connector.Open(ctx, cfg)
}
