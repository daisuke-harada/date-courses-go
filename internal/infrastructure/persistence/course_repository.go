package persistence

import (
	"context"
	"log/slog"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"gorm.io/gorm"
)

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) repository.CourseRepository {
	return &courseRepository{db: db}
}

// visibleToViewer は「公開コース、または viewerID 自身が作成した非公開コース」に絞る GORM スコープです。
// viewerID が 0（未ログイン）のときは user_id が 0 のレコードが存在しないため、公開コースだけが残ります。
func visibleToViewer(viewerID uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 他の条件と AND で結合されるため、OR は必ず括弧で囲む
		return db.Where("(courses.authority = ? OR courses.user_id = ?)", model.CourseAuthorityPublic, viewerID)
	}
}

func (r *courseRepository) Create(ctx context.Context, course *model.Course) error {
	if err := r.db.WithContext(ctx).Create(course).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "courseRepository.Create succeeded", "course_id", course.ID)
	return nil
}

// FindPublicByUserID は指定ユーザーの公開コースを DuringSpots→DateSpot 込みで返します。
// 他人から見えるコース一覧はすべてこちらを使います。
func (r *courseRepository) FindPublicByUserID(ctx context.Context, userID uint) ([]*model.Course, error) {
	return r.findByUserID(ctx, userID, true)
}

// FindAllByUserID は指定ユーザーのコースを非公開も含めて返します。
// 本人がマイページを開いたときだけ使います。
func (r *courseRepository) FindAllByUserID(ctx context.Context, userID uint) ([]*model.Course, error) {
	return r.findByUserID(ctx, userID, false)
}

func (r *courseRepository) findByUserID(ctx context.Context, userID uint, publicOnly bool) ([]*model.Course, error) {
	var courses []*model.Course
	db := r.db.WithContext(ctx).
		Where("courses.user_id = ?", userID).
		Preload("User").
		Preload("DuringSpots.DateSpot")

	if publicOnly {
		db = db.Where("courses.authority = ?", model.CourseAuthorityPublic)
	}

	if err := db.Find(&courses).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.findByUserID failed", "err", err)
		return nil, err
	}
	return courses, nil
}

// Search はフィルタ条件に基づいてコース一覧を返します。
// デートコース一覧には公開コースだけを載せます。作成者本人であっても
// 自分の非公開コースはここには出さず、マイページからのみ辿れるようにしています。
func (r *courseRepository) Search(ctx context.Context, params repository.CourseSearchParams) ([]*model.Course, error) {
	var courses []*model.Course
	db := r.db.WithContext(ctx).
		Where("courses.authority = ?", model.CourseAuthorityPublic).
		Preload("User").
		Preload("DuringSpots.DateSpot")

	if params.PrefectureID != nil {
		db = db.Joins("JOIN during_spots ON during_spots.course_id = courses.id").
			Joins("JOIN date_spots ON date_spots.id = during_spots.date_spot_id").
			Where("date_spots.prefecture_id = ?", *params.PrefectureID).
			Distinct("courses.*")
	}

	if err := db.Find(&courses).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.Search failed", "err", err)
		return nil, apperror.InternalServerError(err)
	}
	return courses, nil
}

// FindByID は指定IDのコースを返します。
// 他人の非公開コースは存在を隠すため、見つからなかった場合と同じ扱いになります。
func (r *courseRepository) FindByID(ctx context.Context, id, viewerID uint) (*model.Course, error) {
	var course model.Course
	if err := r.db.WithContext(ctx).
		Scopes(visibleToViewer(viewerID)).
		Preload("User").
		Preload("DuringSpots.DateSpot").
		First(&course, id).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.FindByID failed", "err", err)
		return nil, err
	}
	return &course, nil
}

// DeleteByID は指定IDのコースを削除します。
func (r *courseRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Course{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.DeleteByID failed", "err", err)
		return err
	}
	return nil
}
