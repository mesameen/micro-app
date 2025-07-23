package grpc

import (
	"context"

	"github.com/mesameen/micro-app/metadata/pkg/model"
	"github.com/mesameen/micro-app/src/api/gen"
	"github.com/mesameen/micro-app/src/pkg/discovery"
	"github.com/mesameen/micro-app/src/pkg/grpcutil"
	"google.golang.org/grpc/credentials"
)

// Gateway defines a movie metadata gRPC gateway
type Gateway struct {
	registry discovery.Registry
	creds    credentials.TransportCredentials
}

// New creates a new gRPC gateway for a movie metadata service
func New(registry discovery.Registry, creds credentials.TransportCredentials) *Gateway {
	return &Gateway{
		registry: registry,
		creds:    creds,
	}
}

// Get gets movie metadata by movie id
func (g *Gateway) GetMovieDetails(ctx context.Context, id string) (*model.Metadata, error) {
	clientConn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry, g.creds)
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
