package mysql

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mesameen/micro-app/rating/internal/repository"
	"github.com/mesameen/micro-app/rating/pkg/model"
)

// Repository defines a new MySQL based movie metadata repository
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL based repository
func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:root@/movieexample")
	if err != nil {
		return nil, err
	}
	return &Repository{
		db: db,
	}, nil
}

// Get retrieves all ratings for a given record
func (r *Repository) Get(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) ([]*model.Rating, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id, value FROM ratings WHERE record_id=? AND record_type=?", recordID, recordType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	var res []*model.Rating
	for rows.Next() {
		var userID string
		var value int32
		if err := rows.Scan(&userID, &value); err != nil {
			return nil, err
		}
		res = append(res, &model.Rating{
			UserID: model.UserID(userID),
			Value:  model.RatingValue(value),
		})
	}
	return res, nil
}

// Put adds a rating for a given record
func (r *Repository) Put(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	if rating == nil {
		return errors.New("rating is nil")
	}
	_, err := r.db.ExecContext(ctx, "INSERT INTO ratings (record_id, record_type, user_id, value) values (?,?,?,?)",
		recordID, recordType, rating.UserID, rating.Value)
	return err
}
