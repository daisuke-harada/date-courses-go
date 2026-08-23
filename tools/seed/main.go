package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/db"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/hotpepper"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/wikimedia"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/persistence"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
	"gorm.io/gorm"
)

// ─── DateSpot 登録 (HotPepper グルメ API の実データ) ─────────────────────────

// seedPrefectureIDs は seed 対象の都道府県です。
// 全 47 都道府県を回すと API リクエストが増えすぎるため、主要エリアに絞ります。
var seedPrefectureIDs = []int{13, 12, 14, 23, 26, 27, 40, 43}

const (
	// seedSpotsPerCombination は 1 つの「都道府県 × ジャンル」で取得する件数です。
	seedSpotsPerCombination = 3
	// seedMinExistingSpots はこの件数以上すでに存在する組み合わせをスキップする閾値です。
	seedMinExistingSpots = 3
)

// seedDateSpots は HotPepper グルメ API から実在店舗を取得して date_spots に投入します。
// 名前・住所・緯度経度・画像 URL・掲載ページ URL はすべて API のデータをそのまま保存し、
// 画像が無い店舗のみ Wikimedia でフォールバックします。
func seedDateSpots(ctx context.Context, gdb *gorm.DB, recruitAPIKey string) {
	repo := persistence.NewDateSpotRepository(gdb)
	fetcher := external.NewSpotFetcher(
		hotpepper.NewClient(recruitAPIKey),
		wikimedia.NewClient(),
	)
	interactor := usecase.NewBatchCreateDateSpotsInteractor(
		repo,
		fetcher,
		seedMinExistingSpots,
		seedSpotsPerCombination,
	)

	for _, prefectureID := range seedPrefectureIDs {
		pref := master.PrefectureByID(prefectureID)
		if pref == nil {
			slog.WarnContext(ctx, "seedDateSpots: unknown prefecture, skip", "prefecture_id", prefectureID)
			continue
		}

		for _, genre := range master.Genres() {
			input := usecase.BatchCreateDateSpotsInput{
				PrefectureID:   pref.ID,
				PrefectureName: pref.Name,
				PrefCode:       pref.PrefCode,
				GenreID:        genre.ID,
				GenreName:      genre.Name,
			}

			if err := interactor.Execute(ctx, input); err != nil {
				slog.ErrorContext(ctx, "seedDateSpots: execute failed, skipping",
					"prefecture", pref.Name,
					"genre", genre.Name,
					"err", err,
				)
			}
		}
	}
}

// dateSpotIDsByPrefecture は指定した都道府県の date_spot ID を返します (コース seed 用)。
func dateSpotIDsByPrefecture(ctx context.Context, gdb *gorm.DB, prefectureIDs ...int) []uint {
	var ids []uint
	if err := gdb.WithContext(ctx).
		Model(&model.DateSpot{}).
		Where("prefecture_id IN ?", prefectureIDs).
		Order("id").
		Pluck("id", &ids).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotIDsByPrefecture: failed to fetch ids",
			"prefecture_ids", prefectureIDs, "err", err)
		return nil
	}
	return ids
}

// ─── User 登録 ────────────────────────────────────────────────────────────────

type userInput struct {
	Email  string
	Name   string
	Gender model.Gender
	Image  string
}

func seedUsers(ctx context.Context, gdb *gorm.DB) {
	repo := persistence.NewUserRepository(gdb)
	auth := service.NewAuthService()

	// seed 用の共通パスワード（Rails の seed と同じ "foobar"）
	defaultPassword, err := auth.HashPassword("foobar")
	if err != nil {
		slog.ErrorContext(ctx, "seedUsers: failed to hash default password", "err", err)
		return
	}
	// admin 用パスワード（Rails の seed と同じ "adminstrator"）
	adminPassword, err := auth.HashPassword("adminstrator")
	if err != nil {
		slog.ErrorContext(ctx, "seedUsers: failed to hash admin password", "err", err)
		return
	}

	type userSeed struct {
		userInput
		Password string
		Admin    bool
	}

	users := []userSeed{
		{userInput: userInput{Email: "guest@gmail.com", Name: "guest", Gender: "男性", Image: "public/images/user_images/man1.jpg"}, Password: defaultPassword},
		{userInput: userInput{Email: "daisuke@gmail.com", Name: "daisuke", Gender: "男性", Image: "public/images/user_images/man2.jpg"}, Password: defaultPassword},
		{userInput: userInput{Email: "kenta@gmail.com", Name: "peter", Gender: "男性", Image: "public/images/user_images/spiderman.png"}, Password: defaultPassword},
		{userInput: userInput{Email: "marika@gmail.com", Name: "marika", Gender: "女性", Image: "public/images/user_images/woman1.jpg"}, Password: defaultPassword},
		{userInput: userInput{Email: "nanase@gmail.com", Name: "nanase", Gender: "女性", Image: "public/images/user_images/woman2.jpg"}, Password: defaultPassword},
		{userInput: userInput{Email: "kanakana@gmail.com", Name: "kanakana", Gender: "女性", Image: "public/images/user_images/woman3.jpg"}, Password: defaultPassword},
		{userInput: userInput{Email: "adminstrator@gmail.com", Name: "admin", Gender: "男性", Image: "public/images/user_images/man1.jpg"}, Password: adminPassword, Admin: true},
	}

	// test1〜12 (男性)
	for i := 1; i <= 12; i++ {
		users = append(users, userSeed{
			userInput: userInput{
				Email:  fmt.Sprintf("%d@gmail.com", i),
				Name:   fmt.Sprintf("test%d", i),
				Gender: "男性",
				Image:  fmt.Sprintf("public/images/user_images/man%d.jpg", rand.Intn(3)+1),
			},
			Password: defaultPassword,
		})
	}
	// test13〜24 (女性)
	for i := 13; i <= 24; i++ {
		users = append(users, userSeed{
			userInput: userInput{
				Email:  fmt.Sprintf("%d@gmail.com", i),
				Name:   fmt.Sprintf("test%d", i),
				Gender: "女性",
				Image:  fmt.Sprintf("public/images/user_images/woman%d.jpg", rand.Intn(3)+1),
			},
			Password: defaultPassword,
		})
	}

	for _, u := range users {
		var existing model.User
		if err := gdb.WithContext(ctx).Where("email = ?", u.Email).First(&existing).Error; err == nil {
			slog.InfoContext(ctx, "User already exists, skip", "email", u.Email)
			continue
		}
		img := u.Image
		user := model.User{
			Name:           u.Name,
			Email:          u.Email,
			Gender:         u.Gender,
			Image:          &img,
			Admin:          u.Admin,
			PasswordDigest: u.Password,
		}
		if err := repo.Create(ctx, &user); err != nil {
			slog.ErrorContext(ctx, "seedUsers: create failed", "email", u.Email, "err", err)
		} else {
			slog.InfoContext(ctx, "User created", "email", u.Email, "id", user.ID)
		}
	}
}

// ─── Relationship 登録 ────────────────────────────────────────────────────────

func seedRelationships(ctx context.Context, gdb *gorm.DB) {
	repo := persistence.NewRelationshipRepository(gdb)

	for userID := 1; userID <= 30; userID++ {
		var followIDs []uint
		switch userID {
		case 29:
			followIDs = []uint{30, 1}
		case 30:
			followIDs = []uint{1, 2}
		default:
			followIDs = []uint{uint(userID + 1), uint(userID + 2)}
		}

		for _, followID := range followIDs {
			var existing model.Relationship
			if err := gdb.WithContext(ctx).
				Where("user_id = ? AND follow_id = ?", userID, followID).
				First(&existing).Error; err == nil {
				slog.InfoContext(ctx, "Relationship already exists, skip", "user_id", userID, "follow_id", followID)
				continue
			}
			rel := model.Relationship{
				UserID:   uint(userID),
				FollowID: followID,
			}
			if err := repo.Create(ctx, &rel); err != nil {
				slog.ErrorContext(ctx, "seedRelationships: create failed",
					"user_id", userID, "follow_id", followID, "err", err)
			} else {
				slog.InfoContext(ctx, "Relationship created", "user_id", userID, "follow_id", followID)
			}
		}
	}
}

// レビュー機能は無効化中（復活させる場合はコメントアウトを解除する）
// // ─── DateSpotReview 登録 ──────────────────────────────────────────────────────
//
// func seedDateSpotReviews(ctx context.Context, gdb *gorm.DB) {
// 	rates := []float64{2, 2.5, 3, 3.5, 4, 4.5, 5}
// 	// Rails の "test" * 18 相当
// 	content := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest"
//
// 	repo := persistence.NewDateSpotReviewRepository(gdb)
//
// 	// 非管理者ユーザー (adminstrator@gmail.com 以外) を取得
// 	var users []model.User
// 	if err := gdb.WithContext(ctx).
// 		Where("email != ?", "adminstrator@gmail.com").
// 		Find(&users).Error; err != nil {
// 		slog.ErrorContext(ctx, "seedDateSpotReviews: failed to fetch users", "err", err)
// 		return
// 	}
//
// 	// DateSpot の件数を取得
// 	var dateSpotCount int64
// 	gdb.WithContext(ctx).Model(&model.DateSpot{}).Count(&dateSpotCount)
// 	if dateSpotCount == 0 {
// 		slog.Warn("seedDateSpotReviews: no date spots found, skip")
// 		return
// 	}
//
// 	// 各ユーザーにランダムに 5 件レビューを投稿 (重複スキップ)
// 	for _, user := range users {
// 		for i := 0; i < 5; i++ {
// 			// ランダムな DateSpot を取得
// 			var ds model.DateSpot
// 			offset := rand.Int63n(dateSpotCount)
// 			if err := gdb.WithContext(ctx).Offset(int(offset)).First(&ds).Error; err != nil {
// 				continue
// 			}
// 			dateSpotID := ds.ID
//
// 			// 同ユーザー × 同スポットの重複チェック
// 			var existing model.DateSpotReview
// 			if err := gdb.WithContext(ctx).
// 				Where("user_id = ? AND date_spot_id = ?", user.ID, dateSpotID).
// 				First(&existing).Error; err == nil {
// 				continue
// 			}
//
// 			rate := rates[rand.Intn(len(rates))]
// 			review := model.DateSpotReview{
// 				Rate:       &rate,
// 				Content:    &content,
// 				UserID:     user.ID,
// 				DateSpotID: dateSpotID,
// 			}
// 			if err := repo.Create(ctx, &review); err != nil {
// 				slog.ErrorContext(ctx, "seedDateSpotReviews: create failed",
// 					"user_id", user.ID, "date_spot_id", dateSpotID, "err", err)
// 			} else {
// 				slog.InfoContext(ctx, "DateSpotReview created", "user_id", user.ID, "date_spot_id", dateSpotID)
// 			}
// 		}
// 	}
// }

// ─── Course + DuringSpot 登録 ─────────────────────────────────────────────────

func courseCreate(
	ctx context.Context,
	gdb *gorm.DB,
	userIDStart, userIDEnd int,
	dateSpotIDs []uint,
	travelMode string,
) {
	courseRepo := persistence.NewCourseRepository(gdb)
	duringSpotRepo := persistence.NewDuringSpotRepository(gdb)

	var users []model.User
	if err := gdb.WithContext(ctx).
		Where("id BETWEEN ? AND ?", userIDStart, userIDEnd).
		Find(&users).Error; err != nil {
		slog.ErrorContext(ctx, "courseCreate: failed to fetch users", "err", err)
		return
	}

	for _, user := range users {
		// 既にコースを持つユーザーはスキップ
		var existingCourse model.Course
		if err := gdb.WithContext(ctx).
			Where("user_id = ?", user.ID).
			First(&existingCourse).Error; err == nil {
			slog.InfoContext(ctx, "Course already exists, skip", "user_id", user.ID)
			continue
		}

		// dateSpotIDs をシャッフル
		shuffled := make([]uint, len(dateSpotIDs))
		copy(shuffled, dateSpotIDs)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		course := model.Course{
			UserID:     user.ID,
			TravelMode: travelMode,
			Authority:  "公開",
		}
		if err := courseRepo.Create(ctx, &course); err != nil {
			slog.ErrorContext(ctx, "courseCreate: course create failed", "user_id", user.ID, "err", err)
			continue
		}
		slog.InfoContext(ctx, "Course created", "user_id", user.ID, "course_id", course.ID)

		// Rails と同じく index 0, 1, 3 の 3 スポットを DuringSpot に登録
		for _, idx := range []int{0, 1, 3} {
			if idx >= len(shuffled) {
				continue
			}
			ds := model.DuringSpot{
				CourseID:   course.ID,
				DateSpotID: shuffled[idx],
			}
			if err := duringSpotRepo.Create(ctx, &ds); err != nil {
				slog.ErrorContext(ctx, "courseCreate: duringSpot create failed",
					"course_id", course.ID, "date_spot_id", shuffled[idx], "err", err)
			} else {
				slog.InfoContext(ctx, "DuringSpot created",
					"course_id", course.ID, "date_spot_id", shuffled[idx])
			}
		}
	}
}

// ─── main ─────────────────────────────────────────────────────────────────────

func main() {
	logger.Init("date-courses-go-seed", false)
	defer logger.Close()

	ctx := context.Background()
	cfg := config.Get()

	gdb, err := db.Connect(ctx, cfg.DB)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	slog.Info("database connected")

	// ─── DateSpots ────────────────────────────────────────────────────────────
	slog.Info("=== seeding DateSpots (HotPepper) ===")
	seedDateSpots(ctx, gdb, cfg.Recruit.APIKey)

	// ─── Users ───────────────────────────────────────────────────────────────
	slog.Info("=== seeding Users ===")
	seedUsers(ctx, gdb)

	// ─── Relationships ────────────────────────────────────────────────────────
	slog.Info("=== seeding Relationships ===")
	seedRelationships(ctx, gdb)

	// ─── DateSpotReviews ──────────────────────────────────────────────────────
	// レビュー機能は無効化中（復活させる場合はコメントアウトを解除する）
	// slog.Info("=== seeding DateSpotReviews ===")
	// seedDateSpotReviews(ctx, gdb)

	// ─── Courses & DuringSpots ────────────────────────────────────────────────
	slog.Info("=== seeding Courses & DuringSpots ===")
	courseCreate(ctx, gdb, 2, 6, dateSpotIDsByPrefecture(ctx, gdb, 13), "DRIVING")
	courseCreate(ctx, gdb, 4, 6, dateSpotIDsByPrefecture(ctx, gdb, 13, 12), "DRIVING")
	courseCreate(ctx, gdb, 7, 11, dateSpotIDsByPrefecture(ctx, gdb, 27), "BICYCLING")
	courseCreate(ctx, gdb, 9, 11, dateSpotIDsByPrefecture(ctx, gdb, 27, 26), "DRIVING")
	courseCreate(ctx, gdb, 12, 16, dateSpotIDsByPrefecture(ctx, gdb, 27, 26), "DRIVING")
	courseCreate(ctx, gdb, 14, 16, dateSpotIDsByPrefecture(ctx, gdb, 14), "WALKING")
	courseCreate(ctx, gdb, 17, 20, dateSpotIDsByPrefecture(ctx, gdb, 23), "DRIVING")
	courseCreate(ctx, gdb, 21, 25, dateSpotIDsByPrefecture(ctx, gdb, 40), "DRIVING")
	courseCreate(ctx, gdb, 26, 30, dateSpotIDsByPrefecture(ctx, gdb, 40, 43), "DRIVING")

	slog.Info("=== seed completed ===")
}
