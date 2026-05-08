package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/db"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/gemini"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/google_places"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/persistence"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
)

func main() {
	logger.Init("date-courses-go-batch", false)
	defer logger.Close()

	cfg := config.Get()
	ctx := context.Background()

	gormDB, err := db.Connect(ctx, cfg.DB)
	if err != nil {
		slog.Error("batch: failed to connect DB", "err", err)
		os.Exit(1)
	}

	repo := persistence.NewDateSpotRepository(gormDB)
	geminiClient := gemini.NewClient(cfg.Gemini.APIKey, cfg.Gemini.Model)
	placesClient := google_places.NewClient(cfg.GooglePlaces.APIKey)
	fetcher := external.NewSpotFetcher(geminiClient, placesClient)

	interactor := usecase.NewBatchCreateDateSpotsInteractor(
		repo,
		fetcher,
		cfg.Batch.MinExistingSpots,
		cfg.Batch.SpotsPerCombination,
	)

	prefectures := master.Prefectures()
	genres := master.Genres()

	taskCount := 0
	for _, pref := range prefectures {
		for _, genre := range genres {
			if taskCount >= cfg.Batch.MaxTasksPerRun {
				slog.InfoContext(ctx, "batch: reached max tasks per run", "max", cfg.Batch.MaxTasksPerRun)
				return
			}

			input := usecase.BatchCreateDateSpotsInput{
				PrefectureID:   pref.ID,
				PrefectureName: pref.Name,
				GenreID:        genre.ID,
				GenreName:      genre.Name,
			}

			if err := interactor.Execute(ctx, input); err != nil {
				slog.ErrorContext(ctx, "batch: execute failed, skipping",
					"prefecture", pref.Name,
					"genre", genre.Name,
					"err", err,
				)
			}

			taskCount++
			// Nominatim の 1req/sec 制限を守るため少し待機
			time.Sleep(1100 * time.Millisecond)
		}
	}

	slog.InfoContext(ctx, "batch: completed", "tasks", taskCount)
}
