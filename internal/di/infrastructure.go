package di

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/db"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// ProvideDB provides *gorm.DB constructed by infrastructure/db.Connect.
func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	return db.Connect(context.Background(), cfg.DB)
}

// ProvideRepositories は全リポジトリのコンストラクタを Container に登録します。
// dig が *gorm.DB を解決して各リポジトリに注入します。
func ProvideRepositories(ct *Container) {
	ct.MustProvide(persistence.NewUserRepository)
	ct.MustProvide(persistence.NewDateSpotRepository)
	ct.MustProvide(persistence.NewAddressRepository)
	ct.MustProvide(persistence.NewCourseRepository)
	ct.MustProvide(persistence.NewDateSpotReviewRepository)
	ct.MustProvide(persistence.NewDuringSpotRepository)
	ct.MustProvide(persistence.NewRelationshipRepository)
}
