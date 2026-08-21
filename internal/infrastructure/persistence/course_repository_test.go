package persistence_test

import (
	"context"
	"strings"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newDryRunDB は実際には接続しない GORM のセッションと、発行された SQL の記録先を返します。
// DryRun では SQL の組み立てだけが行われるため、発行されるクエリを検証できます。
func newDryRunDB(t *testing.T) (*gorm.DB, *[]string) {
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

	captured := []string{}
	capture := func(d *gorm.DB) {
		captured = append(captured, d.Statement.SQL.String())
	}
	require.NoError(t, db.Callback().Query().After("gorm:query").Register("test:capture_query", capture))
	require.NoError(t, db.Callback().Delete().After("gorm:delete").Register("test:capture_delete", capture))

	return db, &captured
}

// issuedSQL は記録された SQL をまとめて1つの文字列にします。
func issuedSQL(captured *[]string) string {
	return strings.Join(*captured, "\n")
}

func TestCourseRepository_Search(t *testing.T) {
	ctx := context.Background()

	// 一覧には公開コースだけを載せる。作成者本人の非公開コースも出さない
	t.Run("filters_to_public_courses_only", func(t *testing.T) {
		db, captured := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		_, _ = repo.Search(ctx, repository.CourseSearchParams{})

		assert.Contains(t, issuedSQL(captured), "WHERE courses.authority = ?")
		assert.NotContains(t, issuedSQL(captured), "user_id")
	})

	t.Run("keeps_public_filter_with_prefecture_id", func(t *testing.T) {
		db, captured := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)
		prefectureID := 13

		_, _ = repo.Search(ctx, repository.CourseSearchParams{PrefectureID: &prefectureID})

		assert.Contains(t, issuedSQL(captured), "courses.authority = ?")
		assert.Contains(t, issuedSQL(captured), "date_spots.prefecture_id = ?")
	})
}

func TestCourseRepository_FindByID(t *testing.T) {
	ctx := context.Background()

	// 公開 or 自分の非公開。OR は括弧で囲まれていないと AND と結合して条件が壊れる
	t.Run("filters_by_visibility", func(t *testing.T) {
		db, captured := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		_, _ = repo.FindByID(ctx, 1, 7)

		assert.Contains(t, issuedSQL(captured), "(courses.authority = ? OR courses.user_id = ?)")
	})
}

func TestCourseRepository_FindByUserID(t *testing.T) {
	ctx := context.Background()

	// 他人から見えるマイページ。公開コースだけを返す。
	// 一覧で人数分のクエリにならないよう IN 句でまとめて引く
	t.Run("public_only_filters_by_authority_with_in_clause", func(t *testing.T) {
		db, captured := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		_, _ = repo.FindPublicByUserIDs(ctx, []uint{1, 2, 3})

		assert.Contains(t, issuedSQL(captured), "courses.user_id IN ")
		assert.Contains(t, issuedSQL(captured), "courses.authority = ?")
	})

	t.Run("public_by_user_ids_returns_empty_without_ids", func(t *testing.T) {
		db, captured := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		result, err := repo.FindPublicByUserIDs(ctx, nil)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.Empty(t, issuedSQL(captured), "ID が無ければクエリを発行しない")
	})

	// 本人のマイページ。非公開コースも含めるので authority で絞らない
	t.Run("all_does_not_filter_by_authority", func(t *testing.T) {
		db, captured := newDryRunDB(t)
		repo := persistence.NewCourseRepository(db)

		_, _ = repo.FindAllByUserID(ctx, 1)

		assert.Contains(t, issuedSQL(captured), "courses.user_id = ?")
		assert.NotContains(t, issuedSQL(captured), "authority")
	})
}
