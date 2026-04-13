package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/db"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/persistence"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
	"gorm.io/gorm"
)

// ─── Genre 名マップ (ActiveHash 相当) ────────────────────────────────────────
var genreNames = map[int]string{
	1:  "ショッピングモール",
	2:  "飲食店",
	3:  "カフェ",
	4:  "アウトドア",
	5:  "遊園地",
	6:  "水族館",
	7:  "寿司",
	8:  "居酒屋",
	9:  "焼肉",
	10: "バーベキュー",
	11: "ランドマーク",
	12: "公園",
}

// ─── Prefecture 名マップ (ActiveHash 相当) ────────────────────────────────────
var prefectureNames = map[int]string{
	1: "北海道", 2: "青森県", 3: "岩手県", 4: "宮城県", 5: "秋田県",
	6: "山形県", 7: "福島県", 8: "茨城県", 9: "栃木県", 10: "群馬県",
	11: "埼玉県", 12: "千葉県", 13: "東京都", 14: "神奈川県", 15: "新潟県",
	16: "富山県", 17: "石川県", 18: "福井県", 19: "山梨県", 20: "長野県",
	21: "岐阜県", 22: "静岡県", 23: "愛知県", 24: "三重県", 25: "滋賀県",
	26: "京都府", 27: "大阪府", 28: "兵庫県", 29: "奈良県", 30: "和歌山県",
	31: "鳥取県", 32: "島根県", 33: "岡山県", 34: "広島県", 35: "山口県",
	36: "徳島県", 37: "香川県", 38: "愛媛県", 39: "高知県", 40: "福岡県",
	41: "佐賀県", 42: "長崎県", 43: "熊本県", 44: "大分県", 45: "宮崎県",
	46: "鹿児島県", 47: "沖縄県",
}

// ─── Geocoding ───────────────────────────────────────────────────────────────

type geocodeResult struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

// geocode は Google Maps Geocoding API を呼び出して緯度・経度を返します。
// API キーが空または "your_api_key_here" の場合は nil を返します。
func geocode(apiKey, address string) (lat *float64, lng *float64) {
	if apiKey == "" || apiKey == "your_api_key_here" {
		return nil, nil
	}
	endpoint := "https://maps.googleapis.com/maps/api/geocode/json"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Warn("geocode: failed to create request", "err", err)
		return nil, nil
	}
	q := url.Values{}
	q.Set("address", address)
	q.Set("key", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Warn("geocode: request failed", "address", address, "err", err)
		return nil, nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result geocodeResult
	if err := json.Unmarshal(body, &result); err != nil || result.Status != "OK" || len(result.Results) == 0 {
		slog.Warn("geocode: unexpected response", "address", address, "status", result.Status)
		return nil, nil
	}
	latVal := result.Results[0].Geometry.Location.Lat
	lngVal := result.Results[0].Geometry.Location.Lng
	return &latVal, &lngVal
}

// ─── 時刻ヘルパー ──────────────────────────────────────────────────────────────

// normalTime は "HH:MM" を 2000-01-01 HH:MM:SS UTC の *time.Time に変換します。
func normalTime(t string) *time.Time {
	if t == "" {
		return nil
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04", "2000-01-01 "+t, time.UTC)
	if err != nil {
		slog.Warn("normalTime: parse failed", "t", t, "err", err)
		return nil
	}
	return &parsed
}

// midnightTime は深夜0〜5時帯を 2000-01-02 HH:MM:SS UTC に変換します。
func midnightTime(t string) *time.Time {
	if t == "" {
		return nil
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04", "2000-01-02 "+t, time.UTC)
	if err != nil {
		slog.Warn("midnightTime: parse failed", "t", t, "err", err)
		return nil
	}
	return &parsed
}

// ─── DateSpot + Address 登録 ──────────────────────────────────────────────────

func spotAndAddressCreate(
	ctx context.Context,
	gdb *gorm.DB,
	name string,
	genreID int,
	openingTime *time.Time,
	closingTime *time.Time,
	prefectureID int,
	cityName string,
	apiKey string,
) {
	// DateSpot を重複チェックしながら登録 (find_or_create_by 相当)
	var dateSpot model.DateSpot
	if err := gdb.WithContext(ctx).Where("name = ?", name).First(&dateSpot).Error; err != nil {
		// 画像はジャンル名.jpg (ActiveStorage の代替として文字列パスを保存)
		imagePath := fmt.Sprintf("public/images/date_spot_images/%s.jpg", genreNames[genreID])
		gid := genreID
		pid := prefectureID
		fullCityName := prefectureNames[prefectureID] + cityName
		lat, lng := geocode(apiKey, fullCityName)
		dateSpot = model.DateSpot{
			GenreID:      &gid,
			PrefectureID: &pid,
			Name:         name,
			CityName:     fullCityName,
			Image:        &imagePath,
			Latitude:     lat,
			Longitude:    lng,
			OpeningTime:  openingTime,
			ClosingTime:  closingTime,
		}
		repo := persistence.NewDateSpotRepository(gdb)
		if err := repo.Create(ctx, &dateSpot); err != nil {
			slog.ErrorContext(ctx, "spotAndAddressCreate: dateSpot create failed", "name", name, "err", err)
			return
		}
		slog.InfoContext(ctx, "DateSpot created", "name", name, "id", dateSpot.ID)
	} else {
		slog.InfoContext(ctx, "DateSpot already exists, skip", "name", name, "id", dateSpot.ID)
	}
}

// ─── User 登録 ────────────────────────────────────────────────────────────────

type userInput struct {
	Email  string
	Name   string
	Gender string
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

// ─── DateSpotReview 登録 ──────────────────────────────────────────────────────

func seedDateSpotReviews(ctx context.Context, gdb *gorm.DB) {
	rates := []float64{2, 2.5, 3, 3.5, 4, 4.5, 5}
	// Rails の "test" * 18 相当
	content := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest"

	repo := persistence.NewDateSpotReviewRepository(gdb)

	// 非管理者ユーザー (adminstrator@gmail.com 以外) を取得
	var users []model.User
	if err := gdb.WithContext(ctx).
		Where("email != ?", "adminstrator@gmail.com").
		Find(&users).Error; err != nil {
		slog.ErrorContext(ctx, "seedDateSpotReviews: failed to fetch users", "err", err)
		return
	}

	// DateSpot の件数を取得
	var dateSpotCount int64
	gdb.WithContext(ctx).Model(&model.DateSpot{}).Count(&dateSpotCount)
	if dateSpotCount == 0 {
		slog.Warn("seedDateSpotReviews: no date spots found, skip")
		return
	}

	// 各ユーザーにランダムに 5 件レビューを投稿 (重複スキップ)
	for _, user := range users {
		for i := 0; i < 5; i++ {
			// ランダムな DateSpot を取得
			var ds model.DateSpot
			offset := rand.Int63n(dateSpotCount)
			if err := gdb.WithContext(ctx).Offset(int(offset)).First(&ds).Error; err != nil {
				continue
			}
			dateSpotID := ds.ID

			// 同ユーザー × 同スポットの重複チェック
			var existing model.DateSpotReview
			if err := gdb.WithContext(ctx).
				Where("user_id = ? AND date_spot_id = ?", user.ID, dateSpotID).
				First(&existing).Error; err == nil {
				continue
			}

			rate := rates[rand.Intn(len(rates))]
			review := model.DateSpotReview{
				Rate:       &rate,
				Content:    &content,
				UserID:     user.ID,
				DateSpotID: dateSpotID,
			}
			if err := repo.Create(ctx, &review); err != nil {
				slog.ErrorContext(ctx, "seedDateSpotReviews: create failed",
					"user_id", user.ID, "date_spot_id", dateSpotID, "err", err)
			} else {
				slog.InfoContext(ctx, "DateSpotReview created", "user_id", user.ID, "date_spot_id", dateSpotID)
			}
		}
	}
}

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

// rangeIDs は start〜end の連続した uint スライスを返します。
func rangeIDs(start, end int) []uint {
	ids := make([]uint, 0, end-start+1)
	for i := start; i <= end; i++ {
		ids = append(ids, uint(i))
	}
	return ids
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

	apiKey := cfg.GoogleMaps.APIKey

	// ─── DateSpot & Address ───────────────────────────────────────────────────
	slog.Info("=== seeding DateSpots & Addresses ===")

	// 東京 1〜6
	spotAndAddressCreate(ctx, gdb, "東京スカイツリー", 11, normalTime("10:00"), normalTime("21:00"), 13, "墨田区押上１丁目１−２", apiKey)
	spotAndAddressCreate(ctx, gdb, "恵比寿ガーデンプレイス", 1, normalTime("07:00"), midnightTime("00:00"), 13, "渋谷区恵比寿4丁目20 ガーデンプレイス", apiKey)
	spotAndAddressCreate(ctx, gdb, "プレゴ・プレゴ", 2, normalTime("16:00"), normalTime("23:00"), 13, "新宿区新宿3-31-3 ＮＳプラザ中央　4F", apiKey)
	spotAndAddressCreate(ctx, gdb, "酒場シナトラ 東京駅店", 9, normalTime("11:00"), normalTime("23:00"), 13, "千代田区丸の内１丁目９−１ 東京駅一番街 2階", apiKey)
	spotAndAddressCreate(ctx, gdb, "カフェ バッハ", 3, normalTime("10:00"), normalTime("19:00"), 13, "台東区日本堤１丁目２３−９", apiKey)
	spotAndAddressCreate(ctx, gdb, "おもてなしとりよし 西新宿店", 2, normalTime("17:00"), normalTime("23:00"), 13, "新宿区西新宿1-10-2 110ビル 11F", apiKey)

	// 千葉 7
	spotAndAddressCreate(ctx, gdb, "東京ディズニーランド", 5, normalTime("10:00"), normalTime("19:00"), 12, "浦安市舞浜１−１", apiKey)

	// 大阪 8〜13
	spotAndAddressCreate(ctx, gdb, "純喫茶 アメリカン", 3, normalTime("09:00"), normalTime("22:00"), 27, "大阪市中央区道頓堀１丁目７−４", apiKey)
	spotAndAddressCreate(ctx, gdb, "ユニバーサル・スタジオ・ジャパン", 5, normalTime("11:00"), normalTime("19:00"), 27, "大阪市此花区桜島２丁目１−３３", apiKey)
	spotAndAddressCreate(ctx, gdb, "焼肉Lab 梅田店", 9, normalTime("12:00"), normalTime("23:00"), 27, "大阪市曽根崎2-10-21 第3河合ビル3F", apiKey)
	spotAndAddressCreate(ctx, gdb, "居酒屋 牡蠣 やまと", 8, normalTime("16:00"), normalTime("23:00"), 27, "大阪市阿倍野区旭町2-1-2 あべのポンテ1F", apiKey)
	spotAndAddressCreate(ctx, gdb, "創蔵", 7, normalTime("17:00"), normalTime("22:00"), 27, "大阪市中央区難波4-6-10", apiKey)
	spotAndAddressCreate(ctx, gdb, "りんくうプレミアム・アウトレット", 1, normalTime("11:00"), normalTime("21:00"), 27, "泉佐野市りんくう往来南３−２８", apiKey)

	// 京都 14〜19
	spotAndAddressCreate(ctx, gdb, "京都タワー", 11, normalTime("09:00"), normalTime("21:00"), 26, "京都市下京区烏丸通七条下る 東塩小路町 721-1", apiKey)
	spotAndAddressCreate(ctx, gdb, "ウメ子の家 四条河原町店", 2, normalTime("17:00"), normalTime("23:00"), 26, "京都市下京区四条小橋東入橋本町105 PONTOビル2F", apiKey)
	spotAndAddressCreate(ctx, gdb, "CINQUE IKARIYA（チンクエイカリヤ）", 2, normalTime("17:00"), normalTime("22:00"), 26, "京都市中京区突抜町138-3", apiKey)
	spotAndAddressCreate(ctx, gdb, "京都 焼き鳥 一", 2, normalTime("17:00"), normalTime("22:00"), 26, "京都市中京区四条室町菊水鉾町585 1F", apiKey)
	spotAndAddressCreate(ctx, gdb, "京都円山　天正", 9, normalTime("17:00"), normalTime("22:00"), 26, "京都市東山区祇園町北側338", apiKey)
	spotAndAddressCreate(ctx, gdb, "Walden Woods Kyoto", 3, normalTime("08:00"), normalTime("19:00"), 26, "京都市下京区栄町５０８−１", apiKey)

	// 神奈川 20〜25
	spotAndAddressCreate(ctx, gdb, "横浜ランドマークタワー", 11, nil, nil, 14, "横浜市西区みなとみらい２丁目２−１", apiKey)
	spotAndAddressCreate(ctx, gdb, "横浜・八景島シーパラダイス", 6, normalTime("11:00"), normalTime("17:00"), 14, "横浜市金沢区八景島", apiKey)
	spotAndAddressCreate(ctx, gdb, "よこはまコスモワールド", 5, normalTime("11:00"), normalTime("21:00"), 14, "横浜市中区新港２丁目８−１", apiKey)
	spotAndAddressCreate(ctx, gdb, "海の公園 バーベキュー場", 10, normalTime("10:00"), normalTime("19:00"), 14, "横浜市金沢区海の公園１０", apiKey)
	spotAndAddressCreate(ctx, gdb, "新横浜ラーメン博物館", 2, normalTime("11:00"), normalTime("21:00"), 14, "横浜市港北区新横浜２丁目１４−２１", apiKey)
	spotAndAddressCreate(ctx, gdb, "旅情個室空間 酒の友 新横浜店", 8, normalTime("17:00"), normalTime("21:00"), 14, "横浜市港北区新横浜3-17-15 3F", apiKey)

	// 愛知 26〜31
	spotAndAddressCreate(ctx, gdb, "名古屋港水族館", 6, normalTime("09:30"), normalTime("17:30"), 23, "名古屋市港区港町1-3", apiKey)
	spotAndAddressCreate(ctx, gdb, "博物館 明治村", 2, normalTime("09:30"), normalTime("17:00"), 23, "犬山市内山1", apiKey)
	spotAndAddressCreate(ctx, gdb, "茶寮 花の宴", 2, normalTime("11:30"), normalTime("22:00"), 23, "安城市大東町１７−８", apiKey)
	spotAndAddressCreate(ctx, gdb, "食堂うさぎや", 2, normalTime("11:00"), normalTime("14:00"), 23, "名古屋市名東区高社２丁目９７", apiKey)
	spotAndAddressCreate(ctx, gdb, "THE ONE AND ONLY", 8, normalTime("18:00"), midnightTime("01:00"), 23, "名古屋市西区牛島町６−１ 名古屋ルーセントタワ 40F", apiKey)
	spotAndAddressCreate(ctx, gdb, "完全個室ダイニング カーヴ隠れや 名古屋駅店", 8, normalTime("17:30"), midnightTime("02:00"), 23, "名古屋市中村区名駅３丁目１５−１１ Ｍ三ダイニングビル 4F", apiKey)

	// 福岡 32〜37
	spotAndAddressCreate(ctx, gdb, "キャナルシティ博多", 1, normalTime("08:00"), normalTime("23:00"), 40, "福岡市博多区住吉1丁目2", apiKey)
	spotAndAddressCreate(ctx, gdb, "つなぐダイニング ZINO 天神店", 8, normalTime("20:00"), midnightTime("05:00"), 40, "福岡市中央区大名1-11-22-1", apiKey)
	spotAndAddressCreate(ctx, gdb, "大濠公園", 12, nil, nil, 40, "福岡市中央区大濠公園", apiKey)
	spotAndAddressCreate(ctx, gdb, "マリンワールド海の中道", 6, nil, nil, 40, "福岡市東区西戸崎18-28", apiKey)
	spotAndAddressCreate(ctx, gdb, "麺劇場 玄瑛", 2, normalTime("11:30"), normalTime("21:00"), 40, "福岡市中央区薬院 2-16-3", apiKey)
	spotAndAddressCreate(ctx, gdb, "芥屋の大門", 4, nil, nil, 40, "糸島市志摩芥屋６７５−２", apiKey)

	// 熊本 38
	spotAndAddressCreate(ctx, gdb, "あか牛丼いわさき", 2, normalTime("11:00"), midnightTime("00:00"), 43, "阿蘇市乙姫2006-2", apiKey)

	// ─── Users ───────────────────────────────────────────────────────────────
	slog.Info("=== seeding Users ===")
	seedUsers(ctx, gdb)

	// ─── Relationships ────────────────────────────────────────────────────────
	slog.Info("=== seeding Relationships ===")
	seedRelationships(ctx, gdb)

	// ─── DateSpotReviews ──────────────────────────────────────────────────────
	slog.Info("=== seeding DateSpotReviews ===")
	seedDateSpotReviews(ctx, gdb)

	// ─── Courses & DuringSpots ────────────────────────────────────────────────
	slog.Info("=== seeding Courses & DuringSpots ===")
	courseCreate(ctx, gdb, 2, 6, rangeIDs(1, 6), "DRIVING")
	courseCreate(ctx, gdb, 4, 6, rangeIDs(1, 7), "DRIVING")
	courseCreate(ctx, gdb, 7, 11, rangeIDs(8, 13), "BICYCLING")
	courseCreate(ctx, gdb, 9, 11, rangeIDs(8, 19), "DRIVING")
	courseCreate(ctx, gdb, 12, 16, rangeIDs(8, 19), "DRIVING")
	courseCreate(ctx, gdb, 14, 16, rangeIDs(20, 25), "WALKING")
	courseCreate(ctx, gdb, 17, 20, rangeIDs(26, 31), "DRIVING")
	courseCreate(ctx, gdb, 21, 25, rangeIDs(32, 37), "DRIVING")
	courseCreate(ctx, gdb, 26, 30, rangeIDs(32, 38), "DRIVING")

	slog.Info("=== seed completed ===")
}
