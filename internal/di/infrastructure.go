package di

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/db"
	"gorm.io/gorm"
)

// ProvideDB provides *gorm.DB constructed by infrastructure/db.Connect.
func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	return db.Connect(context.Background(), cfg.DB)
}
