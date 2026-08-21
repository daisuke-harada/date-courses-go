package persistence

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteDateSpot(t *testing.T) {
	// date_spots を参照している子レコードを先に消さないと外部キー制約で削除できない
	t.Run("deletes_children_before_date_spot", func(t *testing.T) {
		db, captured := newDryRunDBForDelete(t)

		_ = deleteDateSpot(db, 3)

		require.Equal(t, 3, len(*captured), "レビュー・コース中間テーブル・本体で3回の DELETE が必要")

		sqls := *captured
		assert.Contains(t, sqls[0], "DELETE FROM `date_spot_reviews`")
		assert.Contains(t, sqls[1], "DELETE FROM `during_spots`")
		assert.Contains(t, sqls[2], "DELETE FROM `date_spots`")
	})

	t.Run("deletes_in_dependency_order", func(t *testing.T) {
		db, captured := newDryRunDBForDelete(t)

		_ = deleteDateSpot(db, 3)

		all := strings.Join(*captured, "\n")
		reviews := strings.Index(all, "DELETE FROM `date_spot_reviews`")
		during := strings.Index(all, "DELETE FROM `during_spots`")
		spot := strings.Index(all, "DELETE FROM `date_spots`")

		assert.Less(t, reviews, spot, "レビューはスポットより先")
		assert.Less(t, during, spot, "コース中間テーブルはスポットより先")
	})
}
