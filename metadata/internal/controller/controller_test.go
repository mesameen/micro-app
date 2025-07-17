package controller

import (
	"errors"
	"testing"

	"github.com/mesameen/micro-app/metadata/internal/repository"
	"github.com/mesameen/micro-app/metadata/pkg/model"
	gen "github.com/mesameen/micro-app/src/api/gen/mock/metadata/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController(t *testing.T) {
	tests := []struct {
		name       string
		expRepoRes *model.Metadata
		expRepoErr error
		wantRes    *model.Metadata
		wantErr    error
	}{
		{
			name:       "not found",
			expRepoErr: repository.ErrNotFound,
			wantErr:    ErrNotFound,
		},
		{
			name:       "unexpected error",
			expRepoErr: errors.New("unexpected error"),
			wantErr:    errors.New("unexpected error"),
		},
		{
			name:       "success",
			expRepoRes: &model.Metadata{},
			wantRes:    &model.Metadata{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := gen.NewMockmetadataRepository(ctrl)
			c := New(repoMock, nil)
			repoMock.EXPECT().Get(t.Context(), "id").Return(tt.expRepoRes, tt.expRepoErr)
			res, err := c.Get(t.Context(), "id")
			assert.Equal(t, tt.wantRes, res, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}
