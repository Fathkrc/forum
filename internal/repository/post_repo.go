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
func (r *PostRepo) InsertPostWithCategories(ctx context.Context, post *model.Post, categoryIDs []int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		`INSERT INTO posts (user_id, title, body) VALUES (?, ?, ?)`,
		post.UserID,
		post.Title,
		post.Content,
	)
	if err != nil {
		return err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, categoryID := range categoryIDs {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`,
			postID,
			categoryID,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	post.ID = int(postID)
	return nil
}
func (r *PostRepo) GetPostByID(ctx context.Context, postID int64) (*model.Post, error) {
	query := `
				SELECT p.id, p.user_id, p.title, p.body, p.created_at, u.username
				FROM posts p
				JOIN users u ON p.user_id = u.id
				WHERE p.id = ?
			`
	post := &model.Post{}
	if err := r.db.QueryRowContext(ctx, query, postID).
		Scan(&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.AuthorUsername); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepo) ListAllPosts(ctx context.Context) ([]*model.Post, error) {

	query := `
					SELECT p.id, p.user_id, p.title, p.body, p.created_at, u.username
					FROM posts p
					JOIN users u ON p.user_id = u.id
					ORDER BY p.created_at
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
			&post.CreatedAt,
			&post.AuthorUsername); err != nil {
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
					SELECT p.id, p.user_id, p.title, p.body, p.created_at, u.username
					FROM posts p
					JOIN users u ON p.user_id = u.id
					WHERE p.user_id = ?
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
			&post.CreatedAt,
			&post.AuthorUsername); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
func (r *PostRepo) ListLikedPostsByUser(ctx context.Context, userID int64) ([]*model.Post, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, p.user_id, p.title, p.body, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_like pl ON pl.post_id = p.id
		WHERE pl.user_id = ? AND pl.vote = 1
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post

	for rows.Next() {
		post := &model.Post{}

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.AuthorUsername,
		)
		if err != nil {
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
					SELECT p.id, p.user_id, p.title, p.body, p.created_at, u.username
					FROM posts p
					JOIN users u ON p.user_id = u.id
					JOIN post_category pc
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
			&post.CreatedAt,
			&post.AuthorUsername); err != nil {
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
	pattern := "%" + searchTerm + "%"

	rows, err := r.db.QueryContext(ctx, `
			SELECT p.id, p.user_id, p.title, p.body, p.created_at, u.username
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.title LIKE ? OR p.body LIKE ?
			ORDER BY p.created_at DESC
		`, pattern, pattern)
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
			&post.CreatedAt,
			&post.AuthorUsername,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil

}
