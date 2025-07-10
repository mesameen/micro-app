package grpc

import (
	"context"
	"errors"

	"github.com/mesameen/micro-app/rating/internal/controller/rating"
	"github.com/mesameen/micro-app/rating/internal/repository"
	"github.com/mesameen/micro-app/rating/pkg/model"
	"github.com/mesameen/micro-app/src/api/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a gRPC rating API handler
type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

// New creates a new rating gRPC handler
func New(ctrl *rating.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

// GetAggregatedRating returns the aggregated rating for a record
func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Error(codes.InvalidArgument, "req is empty or empty id/type")
	}
	rating, err := h.ctrl.GetAggregatedRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.GetAggregatedRatingResponse{
		RatingValue: rating,
	}, nil
}

// PutRating writes a rating for a given record
func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "req is empty or userID or record is empty")
	}
	if err := h.ctrl.PutRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), &model.Rating{
		UserID: model.UserID(req.RecordId),
		Value:  model.RatingValue(req.RatingValue),
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.PutRatingResponse{}, nil
}
