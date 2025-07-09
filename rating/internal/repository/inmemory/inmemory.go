package inmemory

import (
	"context"

	"github.com/mesameen/micro-app/rating/internal/repository"
	"github.com/mesameen/micro-app/rating/pkg/model"
)

// Repository defines a rating repository
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]*model.Rating
}

// New creates a new in memory repository
func New() *Repository {
	return &Repository{
		data: make(map[model.RecordType]map[model.RecordID][]*model.Rating),
	}
}

// Get retrieves all ratings for a given record
func (r *Repository) Get(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) ([]*model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	var ratings []*model.Rating
	var ok bool
	if ratings, ok = r.data[recordType][recordID]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return ratings, nil
}

// Put adds a rating for a given record
func (r *Repository) Put(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	// if recordType isn't exists create new and add recordid with default ratings
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]*model.Rating{}
	}
	r.data[recordType][recordID] = append(r.data[recordType][recordID], rating)
	return nil
}
