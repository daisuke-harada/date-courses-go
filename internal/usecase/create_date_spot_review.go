package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type CreateDateSpotReviewInputPort interface {
	Execute(context.Context, CreateDateSpotReviewInput) (*CreateDateSpotReviewOutput, error)
}

type CreateDateSpotReviewInput struct {
	UserID     uint
	DateSpotID uint
	Rate       *float64
	Content    *string
}

func (i *CreateDateSpotReviewInput) Validate() error {
	var errs []string
	if i.UserID == 0 {
		errs = append(errs, "ユーザーIDを入力してください")
	}
	if i.DateSpotID == 0 {
		errs = append(errs, "デートスポットIDを入力してください")
	}
	if i.Rate != nil {
		if *i.Rate < 0 || *i.Rate > 5 {
			errs = append(errs, "rate は 0 以上 5 以下で指定してください")
		}
	}
	if i.Content != nil {
		if strings.TrimSpace(*i.Content) == "" {
			errs = append(errs, "content を入力してください")
		} else if len(*i.Content) > 1000 {
			errs = append(errs, "content は1000文字以内で入力してください")
		}
	}
	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}
	return nil
}

// NewCreateDateSpotReviewInputFromStrings builds CreateDateSpotReviewInput from raw string values.
// It performs type parsing (user_id/date_spot_id/rate) and returns BadRequest on parse errors.
func NewCreateDateSpotReviewInputFromStrings(userIDStr, dateSpotIDStr, rateStr, contentStr string) (CreateDateSpotReviewInput, error) {
	// parse user id
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return CreateDateSpotReviewInput{}, apperror.BadRequest("user_id は整数で指定してください")
	}
	// parse date spot id
	dateSpotID, err := strconv.Atoi(dateSpotIDStr)
	if err != nil {
		return CreateDateSpotReviewInput{}, apperror.BadRequest("date_spot_id は整数で指定してください")
	}

	// parse rate (optional)
	var rate *float64
	if rateStr != "" {
		r, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			return CreateDateSpotReviewInput{}, apperror.BadRequest("rate は数値で指定してください")
		}
		rate = &r
	}

	// content (optional)
	var content *string
	if contentStr != "" {
		s := contentStr
		content = &s
	}

	return CreateDateSpotReviewInput{
		UserID:     uint(userID),
		DateSpotID: uint(dateSpotID),
		Rate:       rate,
		Content:    content,
	}, nil
}

type CreateDateSpotReviewOutput struct {
	ReviewID        uint
	DateSpotReviews []*model.DateSpotReview
}

type CreateDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewCreateDateSpotReviewUsecase(dateSpotReviewRepository repository.DateSpotReviewRepository) CreateDateSpotReviewInputPort {
	return &CreateDateSpotReviewInteractor{DateSpotReviewRepository: dateSpotReviewRepository}
}

func (i *CreateDateSpotReviewInteractor) Execute(ctx context.Context, input CreateDateSpotReviewInput) (*CreateDateSpotReviewOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	review := &model.DateSpotReview{
		UserID:     input.UserID,
		DateSpotID: input.DateSpotID,
		Rate:       input.Rate,
		Content:    input.Content,
	}
	if err := i.DateSpotReviewRepository.Create(ctx, review); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	reviews, err := i.DateSpotReviewRepository.FindByDateSpotID(ctx, input.DateSpotID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &CreateDateSpotReviewOutput{
		ReviewID:        review.ID,
		DateSpotReviews: reviews,
	}, nil
}
