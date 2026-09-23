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

func (r *VoteRepo) VotePost(ctx context.Context, userID int64, postID int64, vote int) error {

	query := `
					INSERT INTO post_like (user_id, post_id, vote)
					VALUES (?, ?, ?)
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
) error

func (r *VoteRepo) GetPostVote(
	ctx context.Context,
	userID int64,
	postID int64,
) (int, error)

func (r *VoteRepo) GetCommentVote(
	ctx context.Context,
	userID int64,
	commentID int64,
) (int, error)

func (r *VoteRepo) DeletePostVote(
	ctx context.Context,
	userID int64,
	postID int64,
) error

func (r *VoteRepo) DeleteCommentVote(
	ctx context.Context,
	userID int64,
	commentID int64,
) error

func GetPostVoteCounts(ctx context.Context)    {}
func GetCommentVoteCounts(ctx context.Context) {}
