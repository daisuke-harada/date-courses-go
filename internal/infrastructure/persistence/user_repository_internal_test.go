package persistence

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteUser(t *testing.T) {
	// users を参照している子レコードを先に消さないと外部キー制約で削除できない。
	// during_spots はコース経由の孫レコードなので、コースより先に消す必要がある。
	t.Run("deletes_children_before_user", func(t *testing.T) {
		db, captured := newDryRunDBForDelete(t)

		_ = deleteUser(db, 7)

		require.Equal(t, 5, len(*captured), "孫・子・本体で5回の DELETE が必要")

		sqls := *captured
		assert.Contains(t, sqls[0], "DELETE FROM `during_spots`")
		assert.Contains(t, sqls[0], "SELECT id FROM courses WHERE user_id = ?", "コース経由で孫を特定する")
		assert.Contains(t, sqls[1], "DELETE FROM `courses`")
		assert.Contains(t, sqls[2], "DELETE FROM `date_spot_reviews`")
		assert.Contains(t, sqls[3], "DELETE FROM `relationships`")
		assert.Contains(t, sqls[3], "follow_id = ?", "フォロー・フォロワーの両方を消す")
		assert.Contains(t, sqls[4], "DELETE FROM `users`")
	})

	// 順序が崩れると外部キー制約で失敗するため、並び自体を検証する
	t.Run("deletes_in_dependency_order", func(t *testing.T) {
		db, captured := newDryRunDBForDelete(t)

		_ = deleteUser(db, 7)

		all := strings.Join(*captured, "\n")
		during := strings.Index(all, "DELETE FROM `during_spots`")
		courses := strings.Index(all, "DELETE FROM `courses`")
		users := strings.Index(all, "DELETE FROM `users`")

		assert.Less(t, during, courses, "during_spots はコースより先")
		assert.Less(t, courses, users, "コースはユーザーより先")
	})
}
