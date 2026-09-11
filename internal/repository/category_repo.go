package repository

import (
	"cmd/forum/internal/model"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) InsertCategory(ctx context.Context, category *model.Category) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
        INSERT INTO categories (name)
        VALUES (?)
    `
	result, err := tx.ExecContext(ctx, query,
		strings.TrimSpace(category.Name),
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}
	category.ID = int(id)

	return tx.Commit()

}

func (r *CategoryRepo) InsertMultiCategory(ctx context.Context, categories []*model.Category) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
        INSERT INTO categories (name)
        VALUES (?)
    `
	for _, category := range categories {
		result, err := tx.ExecContext(ctx, query,
			strings.TrimSpace(category.Name),
		)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()

		if err != nil {
			return err
		}
		category.ID = int(id)
	}

	return tx.Commit()
}

func (r *CategoryRepo) ListMultiCategory(ctx context.Context) ([]*model.Category, error) {

	query := ` 
				SELECT id, name
    			FROM categories
    			ORDER BY name
			`
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*model.Category
	for rows.Next() {
		category := &model.Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	categoryID := int64(id)
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `
		DELETE FROM categories
		WHERE id = ?
	`, categoryID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No Category found")
	}

	return tx.Commit()

}

func (r *CategoryRepo) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	query := ` 
				SELECT id, name
    			FROM categories
    			WHERE id = ?
			`
	var category *model.Category

	if err := r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name); err != nil {
		return nil, err
	}

	return category, nil
}
