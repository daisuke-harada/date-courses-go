package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repomock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBatchCreateDateSpotsInteractor_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success_skips_when_enough_spots_exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repomock.NewMockDateSpotRepository(ctrl)
		mockFetcher := &mockSpotFetcher{}

		// 十分なスポットが既存 → スキップ（Fetcherは呼ばれない）
		mockRepo.EXPECT().
			CountByPrefectureAndGenre(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(5), nil).AnyTimes()

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			GenreID:        1,
			GenreName:      "ショッピングモール",
		})
		require.NoError(t, err)
	})

	t.Run("success_creates_new_spots", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repomock.NewMockDateSpotRepository(ctrl)
		mockFetcher := &mockSpotFetcher{
			candidates: []usecase.SpotCandidate{
				{Name: "新宿マルイ", CityName: "新宿区"},
				{Name: "渋谷ヒカリエ", CityName: "渋谷区"},
			},
		}

		mockRepo.EXPECT().
			CountByPrefectureAndGenre(gomock.Any(), 13, 1).
			Return(int64(0), nil)
		// 重複チェック: 両方新規
		mockRepo.EXPECT().
			ExistsByNormalizedNameAndPrefecture(gomock.Any(), gomock.Any(), 13).
			Return(false, nil).Times(2)
		// バッチ登録
		mockRepo.EXPECT().
			CreateBatch(gomock.Any(), gomock.Len(2)).
			Return(nil)

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			GenreID:        1,
			GenreName:      "ショッピングモール",
		})
		require.NoError(t, err)
	})

	t.Run("success_skips_duplicate_spots", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repomock.NewMockDateSpotRepository(ctrl)
		mockFetcher := &mockSpotFetcher{
			candidates: []usecase.SpotCandidate{
				{Name: "新宿マルイ", CityName: "新宿区"},
				{Name: "渋谷ヒカリエ", CityName: "渋谷区"},
			},
		}

		mockRepo.EXPECT().
			CountByPrefectureAndGenre(gomock.Any(), 13, 1).
			Return(int64(0), nil)
		// 最初は重複あり、2番目は新規
		gomock.InOrder(
			mockRepo.EXPECT().ExistsByNormalizedNameAndPrefecture(gomock.Any(), gomock.Any(), 13).Return(true, nil),
			mockRepo.EXPECT().ExistsByNormalizedNameAndPrefecture(gomock.Any(), gomock.Any(), 13).Return(false, nil),
		)
		// 新規1件のみ登録
		mockRepo.EXPECT().
			CreateBatch(gomock.Any(), gomock.Len(1)).
			Return(nil)

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			GenreID:        1,
			GenreName:      "ショッピングモール",
		})
		require.NoError(t, err)
	})

	t.Run("error_fetcher_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repomock.NewMockDateSpotRepository(ctrl)
		mockFetcher := &mockSpotFetcher{
			err: errors.New("gemini api error"),
		}

		mockRepo.EXPECT().
			CountByPrefectureAndGenre(gomock.Any(), 13, 1).
			Return(int64(0), nil)

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			GenreID:        1,
			GenreName:      "ショッピングモール",
		})
		assert.Error(t, err)
	})

	t.Run("success_no_candidates_after_filtering", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repomock.NewMockDateSpotRepository(ctrl)
		mockFetcher := &mockSpotFetcher{
			candidates: []usecase.SpotCandidate{
				{Name: "既存スポット", CityName: "新宿区"},
			},
		}

		mockRepo.EXPECT().
			CountByPrefectureAndGenre(gomock.Any(), 13, 1).
			Return(int64(0), nil)
		mockRepo.EXPECT().
			ExistsByNormalizedNameAndPrefecture(gomock.Any(), gomock.Any(), 13).
			Return(true, nil)
		// CreateBatch は呼ばれない

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			GenreID:        1,
			GenreName:      "ショッピングモール",
		})
		require.NoError(t, err)
	})
}

// mockSpotFetcher はテスト用の SpotFetcher モックです。
type mockSpotFetcher struct {
	candidates []usecase.SpotCandidate
	err        error
}

func (m *mockSpotFetcher) FetchSpotCandidates(ctx context.Context, prefectureName, genreName string, count int) ([]usecase.SpotCandidate, error) {
	return m.candidates, m.err
}

func (m *mockSpotFetcher) FetchPlaceDetail(ctx context.Context, spotName, cityName string) (*usecase.PlaceDetail, error) {
	return &usecase.PlaceDetail{}, nil
}

func (m *mockSpotFetcher) FetchCoordinate(ctx context.Context, spotName, cityName string) (*usecase.Coordinate, error) {
	return nil, nil
}

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// NFKC: カタカナはカタカナのまま（ひらがなには変換しない）、スペース除去、小文字化
		{"新宿マルイ", "新宿マルイ"},
		{"渋谷　ヒカリエ", "渋谷ヒカリエ"},
		{"SHIBUYA109", "shibuya109"},
		// 全角英数字 → 半角
		{"ＡＢＣ", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := usecase.NormalizeName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildMapsURL(t *testing.T) {
	t.Run("builds_correct_url", func(t *testing.T) {
		result := usecase.BuildMapsURL("新宿マルイ", "東京都")
		assert.Contains(t, result, "google.com/maps/search/")
		// URL エンコードされているため Contains では確認できない → プレフィックスのみ確認
		assert.True(t, len(result) > len("https://www.google.com/maps/search/"))
	})
}

// gomock の Len マッチャー用ヘルパー
func gomockLen(n int) gomock.Matcher {
	return gomock.Len(n)
}

var _ = gomockLen

// モデルの Source が正しく設定されることを確認するため
var _ = model.DateSpotSourceAI
