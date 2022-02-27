package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/panta82/goreddit"
)

type PostStore struct {
	*sqlx.DB
}

func (ps *PostStore) Post(id uuid.UUID) (goreddit.Post, error) {
	var post goreddit.Post
	if err := ps.Get(&post, "SELECT * FROM posts WHERE id = $1", id); err != nil {
		return goreddit.Post{}, fmt.Errorf("Failed to get post %s: %w", id.String(), err)
	}
	return post, nil
}

func (ps *PostStore) PostsByThread(threadID uuid.UUID) ([]goreddit.Post, error) {
	var posts []goreddit.Post
	if err := ps.Select(&posts, `
			SELECT *
			FROM posts
			WHERE thread_id = $1`, threadID); err != nil {
		return []goreddit.Post{}, fmt.Errorf("Failed to query posts for thread %s: %w", threadID.String(), err)
	}
	return posts, nil
}

func (ps *PostStore) CreatePost(post *goreddit.Post) error {
	if err := ps.Get(post, `
			INSERT INTO posts (id, title, content, thread_id, votes)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING *`,
		post.ID, post.Title, post.Content, post.ThreadID, post.Votes); err != nil {
		return fmt.Errorf("Failed to create post %s with title \"%s\": %w", post.ID, post.Title, err)
	}
	return nil
}

func (ps *PostStore) UpdatePost(post *goreddit.Post) error {
	if err := ps.Get(post, `
			UPDATE posts
			SET
			    title = $2,
			    content = $3,
			    thread_id = $4,
			    votes = $5
			WHERE id = $1
			RETURNING *`,
		post.ID, post.Title, post.Content, post.ThreadID, post.Votes); err != nil {
		return fmt.Errorf("Failed to update post %s: %w", post.ID, err)
	}
	return nil
}

func (ps *PostStore) DeletePost(id uuid.UUID) error {
	if _, err := ps.Exec("DELETE FROM posts WHERE id = $1", id); err != nil {
		return fmt.Errorf("Failed to delete post %s: %w", id.String(), err)
	}
	return nil
}
