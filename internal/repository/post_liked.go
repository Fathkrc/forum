package repository

import (
	"cmd/forum/internal/model"
	"context"
)

func (r *PostRepo) ListLikedPostsByUser(ctx context.Context, userID int64) ([]*model.Post, error) {
	query := `
				SELECT p.id, p.user_id, p.title, p.body, p.created_at
				FROM posts p JOIN post_like pl ON pl.post_id = p.id
				WHERE pl.user_id = ? AND pl.vote = 1
			`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
