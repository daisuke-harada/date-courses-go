package persistence_test

import (
	"context"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newDryRunDB は実際には接続しない GORM のセッションを返します。
// DryRun では SQL の組み立てだけが行われるため、発行されるクエリを検証できます。
func newDryRunDB(t *testing.T) (*gorm.DB, *string) {
	t.Helper()

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "user:password@tcp(127.0.0.1:3306)/dummy?charset=utf8mb4&parseTime=True",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun: true,
		// 実際には接続しないため、起動時の Ping を無効にする
		DisableAutomaticPing: true,
		Logger:               logger.Discard,
	})
	require.NoError(t, err)

	var captured string
	err = db.Callback().Query().After("gorm:query").Register("test:capture_sql", func(d *gorm.DB) {
		captured = d.Statement.SQL.String()
	})
	require.NoError(t, err)

	return db, &captured
}

func TestCourseRepository_Search(t *testing.T) {
	ctx := context.Background()

	// 一覧には公開コースだけを載せる。作成者本人の非公開コースも出さない
	t.Run("filters_to_public_courses_only", func(t *testing.T) {
		db, sql := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		_, _ = repo.Search(ctx, repository.CourseSearchParams{})

		assert.Contains(t, *sql, "WHERE courses.authority = ?")
		assert.NotContains(t, *sql, "user_id")
	})

	t.Run("keeps_public_filter_with_prefecture_id", func(t *testing.T) {
		db, sql := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)
		prefectureID := 13

		_, _ = repo.Search(ctx, repository.CourseSearchParams{PrefectureID: &prefectureID})

		assert.Contains(t, *sql, "courses.authority = ?")
		assert.Contains(t, *sql, "date_spots.prefecture_id = ?")
	})
}

func TestCourseRepository_FindByID(t *testing.T) {
	ctx := context.Background()

	// 公開 or 自分の非公開。OR は括弧で囲まれていないと AND と結合して条件が壊れる
	t.Run("filters_by_visibility", func(t *testing.T) {
		db, sql := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		_, _ = repo.FindByID(ctx, 1, 7)

		assert.Contains(t, *sql, "(courses.authority = ? OR courses.user_id = ?)")
	})
}

func TestCourseRepository_FindByUserID(t *testing.T) {
	ctx := context.Background()

	// 他人から見えるマイページ。公開コースだけを返す
	t.Run("public_only_filters_by_authority", func(t *testing.T) {
		db, sql := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		_, _ = repo.FindPublicByUserID(ctx, 1)

		assert.Contains(t, *sql, "courses.user_id = ?")
		assert.Contains(t, *sql, "courses.authority = ?")
	})

	// 本人のマイページ。非公開コースも含めるので authority で絞らない
	t.Run("all_does_not_filter_by_authority", func(t *testing.T) {
		db, sql := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		_, _ = repo.FindAllByUserID(ctx, 1)

		assert.Contains(t, *sql, "courses.user_id = ?")
		assert.NotContains(t, *sql, "authority")
	})
}
