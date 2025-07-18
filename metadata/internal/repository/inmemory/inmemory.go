package inmemory

import (
	"context"
	"sync"

	"github.com/mesameen/micro-app/metadata/internal/repository"
	"github.com/mesameen/micro-app/metadata/pkg/model"
)

// Repository defines an in memory movie metadata repository
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates an inmemory repository
func New() *Repository {
	return &Repository{
		data: make(map[string]*model.Metadata),
	}
}

// Get retrieves movie metadata by movie id.
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(_ context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
