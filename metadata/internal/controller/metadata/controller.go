package metadata

import (
	"context"
	"errors"

	"github.com/mesameen/micro-app/metadata/internal/constants"
	"github.com/mesameen/micro-app/metadata/pkg/model"
)

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

// Controller defines a metadata service controller
type Controller struct {
	repo metadataRepository
}

// New creates a metadata service controller
func New(repo metadataRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, constants.ErrNotFound) {
		return nil, constants.ErrNotFound
	}
	return res, err
}
