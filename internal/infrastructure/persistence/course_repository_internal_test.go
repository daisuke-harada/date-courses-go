package persistence

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// deleteCourse は本番ではトランザクションの中で呼ばれるが、DryRun では
// BEGIN に実接続が必要でトランザクションを張れない。そのため発行される SQL の
// 検証は、トランザクションを挟まずこの関数を直接呼んで行う。
func newDryRunDBForDelete(t *testing.T) (*gorm.DB, *[]string) {
	t.Helper()

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "user:password@tcp(127.0.0.1:3306)/dummy?charset=utf8mb4&parseTime=True",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		// GORM は DELETE ごとに既定のトランザクションを張るため、
		// 実接続のない DryRun では無効にしておく
		SkipDefaultTransaction: true,
		Logger:                 logger.Discard,
	})
	require.NoError(t, err)

	captured := []string{}
	require.NoError(t, db.Callback().Delete().After("gorm:delete").Register("test:capture_delete", func(d *gorm.DB) {
		captured = append(captured, d.Statement.SQL.String())
	}))

	return db, &captured
}

func TestDeleteCourse(t *testing.T) {
	// during_spots を先に消さないと外部キー制約で親を削除できない
	t.Run("deletes_during_spots_before_course", func(t *testing.T) {
		db, captured := newDryRunDBForDelete(t)

		_ = deleteCourse(db, 1)

		require.Equal(t, 2, len(*captured), "during_spots とコースの2回の DELETE が必要")
		assert.Contains(t, (*captured)[0], "DELETE FROM `during_spots`")
		assert.Contains(t, (*captured)[0], "course_id = ?")
		assert.Contains(t, (*captured)[1], "DELETE FROM `courses`")
		assert.True(t,
			strings.Index(issued(captured), "during_spots") < strings.Index(issued(captured), "FROM `courses`"),
			"during_spots の削除が先であること")
	})
}

func issued(captured *[]string) string {
	return strings.Join(*captured, "\n")
}
