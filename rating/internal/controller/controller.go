package controller

import (
	"context"
	"errors"

	"github.com/mesameen/micro-app/rating/internal/repository"
	"github.com/mesameen/micro-app/rating/pkg/model"
	"github.com/mesameen/micro-app/src/pkg/logger"
)

// ErrNotFound is returned when no ratings are found for a record
var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(
		ctx context.Context,
		recordID model.RecordID,
		recordType model.RecordType,
	) ([]*model.Rating, error)
	Put(
		ctx context.Context,
		recordID model.RecordID,
		recordType model.RecordType,
		rating *model.Rating,
	) error
}

type ratingIngester interface {
	Ingest(ctx context.Context) (<-chan model.RatingEvent, error)
}

// Controller defines a rating service controller
type Controller struct {
	repo     ratingRepository
	ingester ratingIngester
}

// New creates a rating service controller
func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{
		repo:     repo,
		ingester: ingester,
	}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it.
func (c *Controller) GetAggregatedRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}
	if len(ratings) == 0 {
		return 0.0, nil
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

// PutRating writes a rating for a given record
func (c *Controller) PutRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}

// StartIngestion starts the ingestion of rating events
func (s *Controller) StartIngestion(ctx context.Context) error {
	ch, err := s.ingester.Ingest(ctx)
	if err != nil {
		return err
	}
	for e := range ch {
		logger.Infof("Consumed message: %v", e)
		if err := s.PutRating(ctx, model.RecordID(e.RecordID), model.RecordType(e.RecordType), &model.Rating{
			UserID: e.UserID,
			Value:  e.Value,
		}); err != nil {
			return err
		}
	}
	return nil
}
