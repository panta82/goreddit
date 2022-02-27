package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/panta82/goreddit"
)

type CommentStore struct {
	*sqlx.DB
}

func (cs *CommentStore) Comment(id uuid.UUID) (goreddit.Comment, error) {
	var comment goreddit.Comment
	if err := cs.Get(&comment, "SELECT * FROM comments WHERE id = $1", id); err != nil {
		return goreddit.Comment{}, fmt.Errorf("Failed to get comment %s: %w", id.String(), err)
	}
	return comment, nil
}

func (cs *CommentStore) CommentsByPost(postID uuid.UUID) ([]goreddit.Comment, error) {
	var comments []goreddit.Comment
	if err := cs.Select(&comments, `
			SELECT *
			FROM comments
			WHERE post_id = $1`, postID); err != nil {
		return []goreddit.Comment{}, fmt.Errorf("Failed to query comments for post %s: %w", postID, err)
	}
	return comments, nil
}

func (cs *CommentStore) CreateComment(comment *goreddit.Comment) error {
	if err := cs.Get(&comment, `
			INSERT INTO comments (id, post_id, content, votes)
			VALUES ($1, $2, $3, $4)`,
		comment.ID, comment.PostID, comment.Content, comment.Votes); err != nil {
		return fmt.Errorf("Failed to create comment %s: %w", comment.ID, err)
	}
	return nil
}

func (cs *CommentStore) UpdateComment(comment *goreddit.Comment) error {
	if err := cs.Get(&comment, `
			UPDATE comments
			SET
			    content = $2,
			    votes = $3
			WHERE id = $1`,
		comment.ID, comment.Content, comment.Votes); err != nil {
		return fmt.Errorf("Failed to update comment %s: %w", comment.ID, err)
	}
	return nil
}

func (cs *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := cs.Exec("DELETE FROM comments WHERE id = $1", id); err != nil {
		return fmt.Errorf("Failed to delete comment %s: %w", id.String(), err)
	}
	return nil
}
