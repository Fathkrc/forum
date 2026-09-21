package repository

import "context"

func (r *CommentRepo) UpsertCommentVote(ctx context.Context, userID, commentID int64, vote int) error {
	query := `
				INSERT INTO comment_like (user_id, comment_id, vote)
				VALUES (?, ?, ?)
				ON CONFLICT(user_id, comment_id) DO UPDATE SET vote = excluded.vote
			`
	_, err := r.db.ExecContext(ctx, query, userID, commentID, vote)
	return err
}
