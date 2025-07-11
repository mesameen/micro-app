package grpc

import (
	"context"
	"errors"

	"github.com/mesameen/micro-app/metadata/internal/controller"
	"github.com/mesameen/micro-app/metadata/pkg/model"
	"github.com/mesameen/micro-app/src/api/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie metadata gRPC handler
type Handler struct {
	gen.UnimplementedMetadataServiceServer
	svc *controller.Controller
}

// New creates a new movie metadata gRPC handler
func New(ctrl *controller.Controller) *Handler {
	return &Handler{
		svc: ctrl,
	}
}

// GetMetadata returns movie metadata by id
func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "request is nil or empty movie_id")
	}
	m, err := h.svc.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.GetMetadataResponse{
		Metadata: model.MetadataToProto(m),
	}, nil
}
