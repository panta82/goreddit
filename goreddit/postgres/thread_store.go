package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/panta82/goreddit"
)

type ThreadStore struct {
	*sqlx.DB
}

func (ts *ThreadStore) Thread(id uuid.UUID) (goreddit.Thread, error) {
	var thread goreddit.Thread
	if err := ts.Get(&thread, "SELECT * FROM threads WHERE id = $1", id); err != nil {
		return goreddit.Thread{}, fmt.Errorf("Failed to get thread %s: %w", id.String(), err)
	}
	return thread, nil
}

func (ts *ThreadStore) Threads() ([]goreddit.Thread, error) {
	var threads []goreddit.Thread
	if err := ts.Select(&threads, "SELECT * FROM threads"); err != nil {
		return []goreddit.Thread{}, fmt.Errorf("Failed to query threads: %w", err)
	}
	return threads, nil
}

func (ts *ThreadStore) CreateThread(thread *goreddit.Thread) error {
	if err := ts.Get(&thread, `
			INSERT INTO threads (id, title, description)
			VALUES ($1, $2, $3)`,
		thread.ID, thread.Title, thread.Description); err != nil {
		return fmt.Errorf("Failed to create thread %s with title \"%s\": %w", thread.ID, thread.Title, err)
	}
	return nil
}

func (ts *ThreadStore) UpdateThread(thread *goreddit.Thread) error {
	if err := ts.Get(&thread, `
			UPDATE threads
			SET
			    title = $2,
			    description = $3
			WHERE id = $1`,
		thread.ID, thread.Title, thread.Description); err != nil {
		return fmt.Errorf("Failed to update thread %s: %w", thread.ID, err)
	}
	return nil
}

func (ts *ThreadStore) DeleteThread(id uuid.UUID) error {
	if _, err := ts.Exec("DELETE FROM threads WHERE id = $1", id); err != nil {
		return fmt.Errorf("Failed to delete thread %s: %w", id.String(), err)
	}
	return nil
}
