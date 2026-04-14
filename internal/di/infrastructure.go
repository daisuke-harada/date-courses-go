package di

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/db"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/persistence"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"gorm.io/gorm"
)

// ProvideDB provides *gorm.DB constructed by infrastructure/db.Connect.
func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	return db.Connect(context.Background(), cfg.DB)
}

// ProvideRepositories は全リポジトリのコンストラクタを Container に登録します。
func ProvideRepositories(ct *Container) {
	ct.MustProvide(persistence.NewUserRepository)
	ct.MustProvide(persistence.NewDateSpotRepository)
	ct.MustProvide(persistence.NewCourseRepository)
	ct.MustProvide(persistence.NewDateSpotReviewRepository)
	ct.MustProvide(persistence.NewDuringSpotRepository)
	ct.MustProvide(persistence.NewRelationshipRepository)
}

// ProvideServices は全ドメインサービスのコンストラクタを Container に登録します。
func ProvideServices(ct *Container) {
	ct.MustProvide(service.NewAuthService)
}

// ProvideUsecases は全ユースケースのコンストラクタを Container に登録します。
func ProvideUsecases(ct *Container) {
	ct.MustProvide(usecase.NewGetDateSpotsUsecase)
	ct.MustProvide(usecase.NewSignupUsecase)
	ct.MustProvide(usecase.NewLoginUsecase)
	ct.MustProvide(usecase.NewGetUsersUsecase)
	ct.MustProvide(usecase.NewGetUserUsecase)
}
