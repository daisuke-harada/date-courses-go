package di

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/db"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
	"gorm.io/gorm"
)

// ProvideDB provides *gorm.DB constructed by infrastructure/db.Connect.
func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	return db.Connect(context.Background(), cfg.DB)
}

// ProvideRepositories は全リポジトリのコンストラクタを Container に登録します。
// dig が *gorm.DB を解決して各リポジトリに注入します。
func ProvideRepositories(ct *Container) {
	ct.MustProvide(repository.NewUserRepository)
	ct.MustProvide(repository.NewDateSpotRepository)
	ct.MustProvide(repository.NewAddressRepository)
	ct.MustProvide(repository.NewCourseRepository)
	ct.MustProvide(repository.NewDateSpotReviewRepository)
	ct.MustProvide(repository.NewDuringSpotRepository)
	ct.MustProvide(repository.NewRelationshipRepository)
}
