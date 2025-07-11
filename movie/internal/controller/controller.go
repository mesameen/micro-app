package controller

import (
	"context"
	"errors"

	metadataModel "github.com/mesameen/micro-app/metadata/pkg/model"
	"github.com/mesameen/micro-app/movie/internal/gateway"
	"github.com/mesameen/micro-app/movie/pkg/model"
	ratingModel "github.com/mesameen/micro-app/rating/pkg/model"
)

// ErrNotFound is returned when the movie metadata is not found
var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(
		ctx context.Context,
		recordID ratingModel.RecordID,
		recordType ratingModel.RecordType,
	) (float64, error)
	PutRating(
		ctx context.Context,
		recordID ratingModel.RecordID,
		recordType ratingModel.RecordType,
		rating *ratingModel.Rating,
	) error
}

type metadataGateway interface {
	GetMovieDetails(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

// Controller defines a movie service controller
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{
		ratingGateway:   ratingGateway,
		metadataGateway: metadataGateway,
	}
}

// Get returns the movie details including the aggregated rating and movie details
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.GetMovieDetails(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details := &model.MovieDetails{
		Metadata: metadata,
	}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(id), ratingModel.Movie)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		// ratings aren't mandatory proceeding further
	} else if err != nil {
		return nil, err
	}
	details.Rating = &rating
	return details, nil
}
