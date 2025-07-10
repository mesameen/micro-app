package grpc

import (
	"context"

	"github.com/mesameen/micro-app/pkg/discovery"
	"github.com/mesameen/micro-app/rating/pkg/model"
	"github.com/mesameen/micro-app/src/api/gen"
	"github.com/mesameen/micro-app/src/pkg/grpcutil"
)

// Gateway defines an gRPC gateway to rating service
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a rating service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	clientConn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return 0, err
	}
	defer clientConn.Close()
	client := gen.NewRatingServiceClient(clientConn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   string(recordID),
		RecordType: string(recordType),
	})
	if err != nil {
		return 0, err
	}
	return resp.RatingValue, nil
}

// PutRating writes a rating
func (g *Gateway) PutRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	return nil
}
