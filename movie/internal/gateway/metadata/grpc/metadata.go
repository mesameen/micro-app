package grpc

import (
	"context"

	"github.com/mesameen/micro-app/metadata/pkg/model"
	"github.com/mesameen/micro-app/src/api/gen"
	"github.com/mesameen/micro-app/src/pkg/discovery"
	"github.com/mesameen/micro-app/src/pkg/grpcutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
	clientConn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry, insecure.NewCredentials())
	if err != nil {
		return nil, err
	}
	defer clientConn.Close()
	client := gen.NewMetadataServiceClient(clientConn)
	var resp *gen.GetMetadataResponse
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		resp, err = client.GetMetadata(ctx, &gen.GetMetadataRequest{
			MovieId: id,
		})
		if err != nil {
			if shouldRetry(err) {
				continue
			}
			return nil, err
		}
		return model.MetadataFromProto(resp.Metadata), nil
	}
	return nil, err
}

func shouldRetry(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == codes.DeadlineExceeded || e.Code() == codes.ResourceExhausted || e.Code() == codes.Unavailable
}
