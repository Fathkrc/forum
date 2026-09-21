package repository

import (
	"cmd/forum/internal/model"
	"context"
	"fmt"
	"strings"
)

func (r *CommentRepo) LoadLikes(ctx context.Context, comments []*model.Comment) error {
	if len(comments) == 0 {
		return nil
	}

	placeholders := make([]string, len(comments))
	args := make([]any, 0, len(comments))
	for i, comment := range comments {
		placeholders[i] = "?"
		args = append(args, comment.ID)
	}

	query := fmt.Sprintf(`
							SELECT c.id,
							(SELECT COUNT(*) FROM comment_like cl WHERE cl.comment_id = c.id AND cl.vote = 1),
							(SELECT COUNT(*) FROM comment_like cl WHERE cl.comment_id = c.id AND cl.vote = -1)
							FROM comments c
							WHERE c.id IN (%s)
						`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	details := make(map[int]*model.Comment, len(comments))
	for rows.Next() {
		comment := &model.Comment{}
		if err := rows.Scan(&comment.ID, &comment.Likes, &comment.Dislikes); err != nil {
			return err
		}
		details[comment.ID] = comment
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, comment := range comments {
		if detail, ok := details[comment.ID]; ok {
			comment.Likes = detail.Likes
			comment.Dislikes = detail.Dislikes
		}
	}
	return nil
}
