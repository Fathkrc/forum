package repository

import (
	"cmd/forum/internal/model"
	"context"
	"database/sql"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) InsertComment(ctx context.Context, comment *model.Comment) error {

	query := `
					INSERT INTO comments (post_id, user_id,content )
					VALUES (?,?,?)
				`

	result, err := r.db.ExecContext(ctx, query, comment.PostId,
		comment.UserID,
		comment.Content)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = int(id)

	return nil

}

func (r *CommentRepo) ListCommentsByPost(ctx context.Context, postID int64) ([]*model.Comment, error) {

	query := `
				SELECT c.id, c.post_id ,c.user_id, c.content, c.created_at, u.username
				FROM comments c JOIN users u
				ON c.user_id = u.id 
				WHERE c.post_id = ?
				ORDER BY c.created_at
			`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*model.Comment
	for rows.Next() {
		comment := &model.Comment{}
		if err := rows.Scan(&comment.ID,
			&comment.PostId,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.AuthorUsername); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepo) GetCommentByID(ctx context.Context, commentID int64) (*model.Comment, error) {
	query := `
				SELECT id, post_id, user_id, content, created_at 
				FROM comments
				where id = ?
			`
	comment := &model.Comment{}
	if err := r.db.QueryRowContext(ctx, query, commentID).
		Scan(&comment.ID,
			&comment.PostId,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt); err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *CommentRepo) GetCommentCount(ctx context.Context, postID int64) (int, error) {
	query := "	SELECT COUNT(*) from comments WHERE post_id = ?"

	var commentCount int
	if err := r.db.QueryRowContext(ctx, query, postID).Scan(&commentCount); err != nil {
		return 0, err
	}
	return commentCount, nil
}
