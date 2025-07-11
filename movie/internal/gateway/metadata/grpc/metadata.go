package grpc

import (
	"context"

	"github.com/mesameen/micro-app/metadata/pkg/model"
	"github.com/mesameen/micro-app/src/api/gen"
	"github.com/mesameen/micro-app/src/pkg/discovery"
	"github.com/mesameen/micro-app/src/pkg/grpcutil"
)

// Gateway defines a movie metadata gRPC gateway
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a movie metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

// Get gets movie metadata by movie id
func (g *Gateway) GetMovieDetails(ctx context.Context, id string) (*model.Metadata, error) {
	clientConn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer clientConn.Close()
	client := gen.NewMetadataServiceClient(clientConn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{
		MovieId: id,
	})
	if err != nil {
		return nil, err
	}
	return model.MetadaFromProto(resp.Metadata), nil
}
