package repository

import (
	"cmd/forum/internal/model"
	"context"
	"fmt"
	"strings"
)

func (r *PostRepo) LoadDetails(ctx context.Context, posts []*model.Post) error {
	if len(posts) == 0 {
		return nil
	}

	placeholders := make([]string, len(posts))
	args := make([]any, 0, len(posts))
	for i, post := range posts {
		placeholders[i] = "?"
		args = append(args, post.ID)
	}

	query := fmt.Sprintf(`
							SELECT p.id, u.username,
							(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id),
							(SELECT COUNT(*) FROM post_like pl WHERE pl.post_id = p.id AND pl.vote = 1),
							(SELECT COUNT(*) FROM post_like pl WHERE pl.post_id = p.id AND pl.vote = -1)
							FROM posts p
							JOIN users u ON u.id = p.user_id
							WHERE p.id IN (%s)
						`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	details := make(map[int]*model.Post, len(posts))
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(&post.ID, &post.AuthorUsername, &post.CommentCount, &post.Likes, &post.Dislikes); err != nil {
			return err
		}
		details[post.ID] = post
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, post := range posts {
		if detail, ok := details[post.ID]; ok {
			post.AuthorUsername = detail.AuthorUsername
			post.CommentCount = detail.CommentCount
			post.Likes = detail.Likes
			post.Dislikes = detail.Dislikes
		}
	}
	return nil
}
