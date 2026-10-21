package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"errors"
	"strings"
)

type PostService struct {
	postRepo    *repository.PostRepo
	voteRepo    *repository.VoteRepo
	commentRepo *repository.CommentRepo
}

func NewPostService(postRepo *repository.PostRepo,
	voteRepo *repository.VoteRepo,
	commentRepo *repository.CommentRepo) *PostService {
	return &PostService{postRepo: postRepo,
		voteRepo:    voteRepo,
		commentRepo: commentRepo}
}

func (s *PostService) CreatePost(ctx context.Context, userID int, title, content string, categoryIDs []int) error {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if userID <= 0 {
		return errors.New("invalid user ID")
	}
	if title == "" {
		return errors.New("post title cannot be empty")
	}
	if content == "" {
		return errors.New("post content cannot be empty")
	}

	seen := make(map[int]bool)
	for _, id := range categoryIDs {
		if id <= 0 {
			return errors.New("invalid category ID")
		}
		if seen[id] {
			return errors.New("duplicate category ID")
		}
		seen[id] = true
	}

	post := &model.Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	return s.postRepo.InsertPostWithCategories(ctx, post, categoryIDs)
}

func (s *PostService) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	if id <= 0 {
		return nil, errors.New("invalid post ID")
	}
	post, err := s.postRepo.GetPostByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	likes, dislikes, err := s.voteRepo.GetPostVoteCounts(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	commentCount, err := s.commentRepo.GetCommentCount(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	post.Likes = likes
	post.Dislikes = dislikes
	post.CommentCount = commentCount
	return post, nil
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	repoPosts, err := s.postRepo.ListAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0, len(repoPosts))

	for _, post := range repoPosts {
		posts = append(posts, *post)
	}
	for i := range posts {
		likes, dislikes, err :=
			s.voteRepo.GetPostVoteCounts(
				ctx, int64(posts[i].ID),
			)
		if err != nil {
			return nil, err
		}
		commentCount, err := s.commentRepo.GetCommentCount(ctx, int64(posts[i].ID))
		if err != nil {
			return nil, err
		}
		posts[i].CommentCount = commentCount
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes
	}

	return posts, nil
}

func (s *PostService) GetPostsByCategory(ctx context.Context, categoryID int) ([]model.Post, error) {
	if categoryID <= 0 {
		return nil, errors.New("invalid category ID")
	}

	repoPosts, err := s.postRepo.ListPostsByCategoryID(ctx, int64(categoryID))
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0, len(repoPosts))
	for _, post := range repoPosts {
		posts = append(posts, *post)
	}
	for i := range posts {
		likes, dislikes, err :=
			s.voteRepo.GetPostVoteCounts(
				ctx, int64(posts[i].ID),
			)
		if err != nil {
			return nil, err
		}
		commentCount, err := s.commentRepo.GetCommentCount(ctx, int64(posts[i].ID))
		if err != nil {
			return nil, err
		}
		posts[i].CommentCount = commentCount
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes
	}
	return posts, nil
}

func (s *PostService) GetPostsByUser(ctx context.Context, userID int) ([]model.Post, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	repoPosts, err := s.postRepo.ListPostsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0, len(repoPosts))
	for _, post := range repoPosts {
		posts = append(posts, *post)
	}
	for i := range posts {
		likes, dislikes, err :=
			s.voteRepo.GetPostVoteCounts(
				ctx, int64(posts[i].ID),
			)
		if err != nil {
			return nil, err
		}
		commentCount, err := s.commentRepo.GetCommentCount(ctx, int64(posts[i].ID))
		if err != nil {
			return nil, err
		}
		posts[i].CommentCount = commentCount
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes
	}

	return posts, nil
}

func (s *PostService) GetLikedPostsByUser(ctx context.Context,
	userID int,
) ([]model.Post, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	repoPosts, err := s.postRepo.ListLikedPostsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0, len(repoPosts))

	for _, post := range repoPosts {
		posts = append(posts, *post)
	}
	for i := range posts {
		likes, dislikes, err :=
			s.voteRepo.GetPostVoteCounts(
				ctx, int64(posts[i].ID),
			)
		if err != nil {
			return nil, err
		}
		commentCount, err := s.commentRepo.GetCommentCount(ctx, int64(posts[i].ID))
		if err != nil {
			return nil, err
		}
		posts[i].CommentCount = commentCount
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes
	}
	return posts, nil
}

func (s *PostService) SearchPosts(ctx context.Context, query string) ([]model.Post, error) {
	query = strings.TrimSpace(query)

	if query == "" {
		return []model.Post{}, nil
	}

	repoPosts, err := s.postRepo.SearchPosts(ctx, query)
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0, len(repoPosts))

	for _, post := range repoPosts {
		posts = append(posts, *post)
	}
	for i := range posts {
		likes, dislikes, err :=
			s.voteRepo.GetPostVoteCounts(
				ctx, int64(posts[i].ID),
			)
		if err != nil {
			return nil, err
		}
		commentCount, err := s.commentRepo.GetCommentCount(ctx, int64(posts[i].ID))
		if err != nil {
			return nil, err
		}
		posts[i].CommentCount = commentCount
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes
	}

	return posts, nil
}
