package service

import (
	"cmd/forum/internal/repository"
	"context"
	"database/sql"
	"errors"
)

type VoteService struct {
	postRepo    *repository.PostRepo
	commentRepo *repository.CommentRepo
	voteRepo    *repository.VoteRepo
}

func NewVoteService(postRepo *repository.PostRepo, voteRepo *repository.VoteRepo, commentRepo *repository.CommentRepo) *VoteService {
	return &VoteService{postRepo: postRepo, commentRepo: commentRepo, voteRepo: voteRepo}
}

func (s *VoteService) VotePost(ctx context.Context,
	userID,
	postID int,
	isLike bool) error {
	if userID <= 0 || postID <= 0 {
		return errors.New("invalid user or post ID")
	}
	vote := -1
	if isLike {
		vote = 1
	}
	currentVote, err := s.voteRepo.GetPostVote(ctx, int64(userID), int64(postID))
	// if the vote is already like and we post like again it should delete
	// we are assuming user liked by missclick and take it back

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if err == nil && currentVote == vote {
		// here is the users vote is like and retrieving data is also like we delete it
		return s.voteRepo.DeletePostVote(ctx, int64(userID), int64(postID))
	}
	return s.voteRepo.VotePost(
		context.Background(),
		int64(userID),
		int64(postID),
		vote,
	)
}

func (s *VoteService) VoteComment(ctx context.Context,
	userID,
	commentID int,
	isLike bool) error {
	if userID <= 0 || commentID <= 0 {
		return errors.New("invalid user or post ID")
	}
	vote := -1
	if isLike {
		vote = 1
	}
	// same for comments if user has a like vote and send like again we delete its voteRepo
	// like reddit or facebook like most of forum websites
	currentVote, err := s.voteRepo.GetCommentVote(ctx, int64(userID), int64(commentID))

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if err == nil && currentVote == vote {
		return s.voteRepo.DeleteCommentVote(ctx, int64(userID), int64(commentID))
	}
	return s.voteRepo.VoteComment(
		ctx,
		int64(userID),
		int64(commentID),
		vote,
	)
}

func (s *VoteService) GetPostVote(ctx context.Context,
	userID int64,
	postID int64,
) (int, error) {
	return s.voteRepo.GetPostVote(ctx, userID, postID)
}

func (s *VoteService) GetCommentVote(ctx context.Context,
	userID int64,
	commentID int64,
) (int, error) {
	return s.voteRepo.GetCommentVote(ctx, userID, commentID)
}

func (s *VoteService) DeletePostVote(ctx context.Context,
	userID int64,
	postID int64,
) error {
	return s.voteRepo.DeletePostVote(ctx, userID, postID)
}

func (s *VoteService) DeleteCommentVote(
	ctx context.Context,
	userID int64,
	commentID int64,
) error {
	return s.voteRepo.DeleteCommentVote(ctx, userID, commentID)
}

func (s *VoteService) GetPostVoteCounts(
	ctx context.Context,
	postID int64) (int, int, error) {
	return s.voteRepo.GetPostVoteCounts(ctx, postID)
}

func (s *VoteService) GetCommentVoteCounts(
	ctx context.Context,
	commentID int64) (int, int, error) {
	return s.voteRepo.GetCommentVoteCounts(ctx, commentID)
}
