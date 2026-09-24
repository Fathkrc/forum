package repository

import (
	"context"
	"database/sql"
)

type VoteRepo struct {
	db *sql.DB
}

func NewVoteRepo(db *sql.DB) *VoteRepo {
	return &VoteRepo{db: db}
}

func (r *VoteRepo) VotePost(
	ctx context.Context,
	userID int64,
	postID int64,
	vote int,
) error {

	query := `
					INSERT INTO post_like (user_id, post_id, vote)
					VALUES (?, ?, ?)
					ON CONFLICT(user_id, post_id)
					DO UPDATE SET vote = excluded.vote
			`
	if _, err := r.db.ExecContext(ctx, query, userID, postID, vote); err != nil {
		return err
	}
	return nil
}

func (r *VoteRepo) VoteComment(
	ctx context.Context,
	userID int64,
	commentID int64,
	vote int,
) error {

	query := `
					INSERT INTO comment_like (user_id, comment_id, vote)
					VALUES (?, ?, ?)
					ON CONFLICT(user_id, comment_id)
					DO UPDATE SET vote = excluded.vote
			`
	if _, err := r.db.ExecContext(ctx, query, userID, commentID, vote); err != nil {
		return err
	}
	return nil

}

func (r *VoteRepo) GetPostVote(
	ctx context.Context,
	userID int64,
	postID int64,
) (int, error) {

	query := `
				SELECT vote
				FROM post_like
				WHERE user_id = ? AND post_id = ?
			`
	var vote int

	if err := r.db.QueryRowContext(ctx, query, userID, postID).Scan(&vote); err != nil {
		return 0, err
	}
	return vote, nil
}

func (r *VoteRepo) GetCommentVote(
	ctx context.Context,
	userID int64,
	commentID int64,
) (int, error) {

	query := `
				SELECT vote
				FROM comment_like
				WHERE user_id = ? AND comment_id = ?
			`
	var vote int

	if err := r.db.QueryRowContext(ctx, query, userID, commentID).Scan(&vote); err != nil {
		return 0, err
	}
	return vote, nil
}

func (r *VoteRepo) DeletePostVote(
	ctx context.Context,
	userID int64,
	postID int64,
) error {

	query := `
				DELETE FROM post_like
				WHERE user_id = ? AND post_id = ?
			`
	result, err := r.db.ExecContext(ctx, query, userID, postID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *VoteRepo) DeleteCommentVote(
	ctx context.Context,
	userID int64,
	commentID int64,
) error {
	query := `
				DELETE FROM comment_like
				WHERE user_id = ? AND comment_id = ?
			`
	result, err := r.db.ExecContext(ctx, query, userID, commentID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *VoteRepo) GetPostVoteCounts(ctx context.Context, postID int64) (int, int, error) {

	query := `
				SELECT
    			COALESCE(SUM(CASE WHEN vote = 1 THEN 1 ELSE 0 END), 0),
    			COALESCE(SUM(CASE WHEN vote = -1 THEN 1 ELSE 0 END), 0)
				FROM post_like
				WHERE post_id = ?
	    `
	var likes int
	var dislikes int

	if err := r.db.QueryRowContext(ctx, query, postID).Scan(&likes, &dislikes); err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil

}
func (r *VoteRepo) GetCommentVoteCounts(ctx context.Context, commentID int64) (int, int, error) {
	query := `
				SELECT
    			COALESCE(SUM(CASE WHEN vote = 1 THEN 1 ELSE 0 END), 0),
    			COALESCE(SUM(CASE WHEN vote = -1 THEN 1 ELSE 0 END), 0)
				FROM comment_like
				WHERE comment_id = ?
	    `
	var likes int
	var dislikes int

	if err := r.db.QueryRowContext(ctx, query, commentID).Scan(&likes, &dislikes); err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}
