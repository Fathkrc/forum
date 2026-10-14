package service_test

import (
	"cmd/forum/internal/database"
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"cmd/forum/internal/service"
	"context"
	"path/filepath"
	"testing"
)

func TestCreatePost(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.NewDataBase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	// Create a real user for the foreign key.
	result, err := db.ExecContext(ctx,
		`INSERT INTO users (username, email, password)
		 VALUES (?, ?, ?)`,
		"testuser", "test@example.com", "test-hash",
	)
	if err != nil {
		t.Fatal(err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create a category.
	result, err = db.ExecContext(ctx,
		`INSERT INTO categories (name) VALUES (?)`,
		"Fantasy",
	)
	if err != nil {
		t.Fatal(err)
	}

	categoryID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewRepository(db)
	svc := service.New(repo)

	err = svc.CreatePost(
		ctx,
		int(userID),
		"My First Post",
		"Testing post creation",
		[]int{int(categoryID)},
	)
	if err != nil {
		t.Fatal(err)
	}

	var postID int
	err = db.QueryRowContext(ctx,
		`SELECT id FROM posts WHERE title = ?`,
		"My First Post",
	).Scan(&postID)
	if err != nil {
		t.Fatal(err)
	}

	var linkedCategoryID int
	err = db.QueryRowContext(ctx,
		`SELECT category_id FROM post_category WHERE post_id = ?`,
		postID,
	).Scan(&linkedCategoryID)
	if err != nil {
		t.Fatal(err)
	}

	if linkedCategoryID != int(categoryID) {
		t.Fatalf("expected category %d, got %d",
			categoryID, linkedCategoryID)
	}

	t.Logf("Created post %d with category %d", postID, linkedCategoryID)
	post, err := svc.GetPostByID(ctx, postID)
	if err != nil {
		t.Fatalf("GetPostByID failed: %v", err)
	}

	if post.ID != postID {
		t.Fatalf("expected post ID %d, got %d", postID, post.ID)
	}

	if post.Title != "My First Post" {
		t.Fatalf("unexpected post title: %s", post.Title)
	}

	t.Logf("Retrieved post: ID=%d, Title=%s", post.ID, post.Title)
	_ = model.Post{} // Not needed; remove model import and this line.
}
