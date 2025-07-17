package controller

import (
	"errors"
	"testing"

	"github.com/mesameen/micro-app/rating/internal/repository"
	"github.com/mesameen/micro-app/rating/pkg/model"
	gen "github.com/mesameen/micro-app/src/api/gen/mock/rating/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetAggregatedRating(t *testing.T) {
	tests := []struct {
		name        string
		wantRes     float64
		wantErr     error
		wantRepoRes []*model.Rating
		wantRepoErr error
	}{
		{
			name:        "not found",
			wantErr:     ErrNotFound,
			wantRepoErr: repository.ErrNotFound,
		},
		{
			name:        "unexpected error",
			wantErr:     errors.New("unexpected error"),
			wantRepoErr: errors.New("unexpected error"),
		},
		{
			name:        "success",
			wantRes:     0.0,
			wantRepoRes: make([]*model.Rating, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := gen.NewMockratingRepository(ctrl)
			c := New(repoMock, nil)
			repoMock.EXPECT().Get(t.Context(), model.RecordID("id"), model.RecordType("movie")).Return(tt.wantRepoRes, tt.wantRepoErr)
			res, err := c.GetAggregatedRating(t.Context(), model.RecordID("id"), model.RecordType("movie"))
			assert.Equal(t, res, tt.wantRes, tt.name)
			assert.Equal(t, err, tt.wantErr, tt.name)
		})
	}
}
