package repository

import (
	"cmd/forum/internal/model"
	"context"
	"database/sql"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) InsertPost(ctx context.Context, post *model.Post) error {
	query := `
				INSERT INTO posts (user_id, title, body)
				VALUES (?, ?, ?)
			`
	result, err := r.db.ExecContext(ctx, query,
		post.UserID,
		post.Title,
		post.Content)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int(id)

	return nil
}

func (r *PostRepo) GetPostByID(ctx context.Context, postID int64) (*model.Post, error) {
	query := `
				SELECT id, user_id, title, body, created_at 
				FROM posts
				WHERE id = ?
			`
	post := &model.Post{}
	if err := r.db.QueryRowContext(ctx, query, postID).
		Scan(&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepo) ListAllPosts(ctx context.Context) ([]*model.Post, error) {

	query := `
					SELECT id, user_id, title, body, created_at
					FROM posts
					ORDER BY created_at DESC
				`
	rows, err := r.db.QueryContext(ctx, query)
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

func (r *PostRepo) ListPostsByUser(ctx context.Context, userID int64) ([]*model.Post, error) {

	query := `
					SELECT id, user_id, title, body, created_at
					FROM posts
					WHERE user_id = ?
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

func (r *PostRepo) ListPostsByCategoryID(ctx context.Context, categoryID int64) ([]*model.Post, error) {

	query := `
					SELECT p.id, p.user_id, p.title, p.body, p.created_at 
					FROM posts p JOIN post_category pc 
					ON p.id = pc.post_id
					WHERE pc.category_id = ? 
				`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
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

func (r *PostRepo) AddCategoryToPost(ctx context.Context, postID int64, categoryID int64) error {

	query := "INSERT INTO post_category (post_id, category_id) VALUES (?, ?)"

	if _, err := r.db.ExecContext(ctx, query, postID, categoryID); err != nil {
		return err
	}
	return nil

}

func (r *PostRepo) SearchPosts(ctx context.Context, searchTerm string) ([]*model.Post, error) {
	query := `
				SELECT id, user_id, title, body, created_at
				FROM posts
				WHERE title LIKE ? OR body LIKE ?
			`
	rows, err := r.db.QueryContext(ctx, query, "%"+searchTerm+"%", "%"+searchTerm+"%")
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
