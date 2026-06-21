package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/db"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/hotpepper"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/jalan"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/wikimedia"
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
	hotpepperClient := hotpepper.NewClient(cfg.Recruit.APIKey)
	jalanClient := jalan.NewClient(cfg.Recruit.APIKey)
	wikimediaClient := wikimedia.NewClient()
	fetcher := external.NewSpotFetcher(hotpepperClient, jalanClient, wikimediaClient)

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
				PrefCode:       pref.PrefCode,
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
		}
	}

	slog.InfoContext(ctx, "batch: completed", "tasks", taskCount)
}
