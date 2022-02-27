package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/panta82/goreddit"
)

type Store struct {
	goreddit.ThreadStore
	goreddit.PostStore
	goreddit.CommentStore
}

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping database: %v", err)
	}

	return &Store{
		ThreadStore:  &ThreadStore{db},
		PostStore:    &PostStore{db},
		CommentStore: &CommentStore{db},
	}, nil
}
