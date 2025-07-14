package controller

import (
	"context"
	"errors"

	"github.com/mesameen/micro-app/metadata/internal/repository"
	"github.com/mesameen/micro-app/metadata/pkg/model"
	"github.com/mesameen/micro-app/src/pkg/logger"
)

// ErrNotFound is returned when a  requested resource not found
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, id string, metadata *model.Metadata) error
}

// Controller defines a metadata service controller
type Controller struct {
	repo  metadataRepository
	cache metadataRepository
}

// New creates a metadata service controller
func New(repo metadataRepository, cache metadataRepository) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
	}
}

// Get returns movie metadata by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	cacheRes, err := c.cache.Get(ctx, id)
	if err == nil {
		logger.Infof("Returning metadata from a cache for %s", id)
		return cacheRes, nil
	}
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err := c.cache.Put(ctx, id, res); err != nil {
		logger.Infof("Error updating a cache: %v", err)
	}
	return res, err
}
