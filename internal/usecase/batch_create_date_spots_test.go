package usecase_test

import (
	"context"
	"errors"
	"testing"

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

		mockRepo.EXPECT().
			CountByPrefectureAndGenre(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(5), nil).AnyTimes()

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			PrefCode:       "13",
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
		mockRepo.EXPECT().
			ExistsByNormalizedNameAndPrefecture(gomock.Any(), gomock.Any(), 13).
			Return(false, nil).Times(2)
		mockRepo.EXPECT().
			CreateBatch(gomock.Any(), gomock.Len(2)).
			Return(nil)

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			PrefCode:       "13",
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
		gomock.InOrder(
			mockRepo.EXPECT().ExistsByNormalizedNameAndPrefecture(gomock.Any(), gomock.Any(), 13).Return(true, nil),
			mockRepo.EXPECT().ExistsByNormalizedNameAndPrefecture(gomock.Any(), gomock.Any(), 13).Return(false, nil),
		)
		mockRepo.EXPECT().
			CreateBatch(gomock.Any(), gomock.Len(1)).
			Return(nil)

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			PrefCode:       "13",
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
			err: errors.New("api error"),
		}

		mockRepo.EXPECT().
			CountByPrefectureAndGenre(gomock.Any(), 13, 1).
			Return(int64(0), nil)

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			PrefCode:       "13",
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

		interactor := usecase.NewBatchCreateDateSpotsInteractor(mockRepo, mockFetcher, 5, 3)
		err := interactor.Execute(ctx, usecase.BatchCreateDateSpotsInput{
			PrefectureID:   13,
			PrefectureName: "東京都",
			PrefCode:       "13",
			GenreID:        1,
			GenreName:      "ショッピングモール",
		})
		require.NoError(t, err)
	})
}

type mockSpotFetcher struct {
	candidates []usecase.SpotCandidate
	err        error
}

func (m *mockSpotFetcher) FetchSpots(_ context.Context, _, _ string, _ int, _ int) ([]usecase.SpotCandidate, error) {
	return m.candidates, m.err
}

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"新宿マルイ", "新宿マルイ"},
		{"渋谷　ヒカリエ", "渋谷ヒカリエ"},
		{"SHIBUYA109", "shibuya109"},
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
		assert.True(t, len(result) > len("https://www.google.com/maps/search/"))
	})
}
