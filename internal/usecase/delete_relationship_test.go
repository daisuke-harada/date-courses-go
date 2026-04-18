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

func TestDeleteRelationshipInteractor_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		relationshipRepo := repomock.NewMockRelationshipRepository(ctrl)
		relationshipRepo.EXPECT().
			DeleteByUserIDs(ctx, uint(1), uint(2)).
			Return(nil)

		interactor := usecase.NewDeleteRelationshipUsecase(relationshipRepo)
		output, err := interactor.Execute(ctx, usecase.DeleteRelationshipInput{
			UserID:   1,
			FollowID: 2,
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
	})

	t.Run("error_repository_delete_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		relationshipRepo := repomock.NewMockRelationshipRepository(ctrl)
		relationshipRepo.EXPECT().
			DeleteByUserIDs(ctx, uint(1), uint(2)).
			Return(errors.New("db error"))

		interactor := usecase.NewDeleteRelationshipUsecase(relationshipRepo)
		output, err := interactor.Execute(ctx, usecase.DeleteRelationshipInput{
			UserID:   1,
			FollowID: 2,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
