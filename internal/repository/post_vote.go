package repository

import "context"

func (r *PostRepo) UpsertPostVote(ctx context.Context, userID, postID int64, vote int) error {
	query := `
				INSERT INTO post_like (user_id, post_id, vote)
				VALUES (?, ?, ?)
				ON CONFLICT(user_id, post_id) DO UPDATE SET vote = excluded.vote
			`
	_, err := r.db.ExecContext(ctx, query, userID, postID, vote)
	return err
}
