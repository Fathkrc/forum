package service_test

import (
	"cmd/forum/internal/database"
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"cmd/forum/internal/service"
	"context"
	"database/sql"
	"errors"
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
	// Use the embedded Fantasy category.
	var categoryID int64
	err = db.QueryRowContext(ctx,
		`SELECT id FROM categories WHERE name = ?`,
		"Fantasy",
	).Scan(&categoryID)
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
	posts, err := svc.GetPostsByCategory(ctx, int(categoryID))
	if err != nil {
		t.Fatalf("GetPostsByCategory failed: %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	if posts[0].ID != postID {
		t.Fatalf("expected post ID %d, got %d", postID, posts[0].ID)
	}

	t.Logf("Category filter returned post: %s", posts[0].Title)
	t.Logf("Retrieved post: ID=%d, Title=%s", post.ID, post.Title)
	_ = model.Post{} // Not needed; remove model import and this line.
}
func TestPostsByCategoryFiltering(t *testing.T) {
	ctx := context.Background()

	// Create a temporary test database
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.NewDataBase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a test user
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

	// Use the embedded Fantasy category.
	var fantasyID int64
	err = db.QueryRowContext(ctx,
		`SELECT id FROM categories WHERE name = ?`,
		"Fantasy",
	).Scan(&fantasyID)
	if err != nil {
		t.Fatal(err)
	}

	// Use the embedded Science Fiction category.
	var scifiID int64
	err = db.QueryRowContext(ctx,
		`SELECT id FROM categories WHERE name = ?`,
		"Science Fiction",
	).Scan(&scifiID)
	if err != nil {
		t.Fatal(err)
	}

	// Initialize the service
	repo := repository.NewRepository(db)
	svc := service.New(repo)

	// Create a Fantasy post
	err = svc.CreatePost(
		ctx,
		int(userID),
		"My Fantasy Post",
		"This post belongs to Fantasy",
		[]int{int(fantasyID)},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Science Fiction post
	err = svc.CreatePost(
		ctx,
		int(userID),
		"My Sci-Fi Post",
		"This post belongs to Science Fiction",
		[]int{int(scifiID)},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Filter by Fantasy
	fantasyPosts, err := svc.GetPostsByCategory(ctx, int(fantasyID))
	if err != nil {
		t.Fatal(err)
	}

	if len(fantasyPosts) != 1 {
		t.Fatalf("expected 1 Fantasy post, got %d", len(fantasyPosts))
	}

	if fantasyPosts[0].Title != "My Fantasy Post" {
		t.Fatalf("unexpected Fantasy post: %s", fantasyPosts[0].Title)
	}

	// Filter by Science Fiction
	scifiPosts, err := svc.GetPostsByCategory(ctx, int(scifiID))
	if err != nil {
		t.Fatal(err)
	}

	if len(scifiPosts) != 1 {
		t.Fatalf("expected 1 Sci-Fi post, got %d", len(scifiPosts))
	}

	if scifiPosts[0].Title != "My Sci-Fi Post" {
		t.Fatalf("unexpected Sci-Fi post: %s", scifiPosts[0].Title)
	}

	t.Log("Fantasy filter returns only Fantasy posts")
	t.Log("Science Fiction filter returns only Science Fiction posts")
	t.Log("Category filtering works correctly!")
}
func TestGetPostsByUser(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.NewDataBase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.New(repo)

	// Create two users.
	result, err := db.ExecContext(ctx,
		`INSERT INTO users (username, email, password)
		 VALUES (?, ?, ?)`,
		"alice", "alice@example.com", "test-hash",
	)
	if err != nil {
		t.Fatal(err)
	}
	aliceID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	result, err = db.ExecContext(ctx,
		`INSERT INTO users (username, email, password)
		 VALUES (?, ?, ?)`,
		"bob", "bob@example.com", "test-hash",
	)
	if err != nil {
		t.Fatal(err)
	}
	bobID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create posts for each user.
	err = svc.CreatePost(ctx, int(aliceID),
		"Alice's Post", "Written by Alice", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = svc.CreatePost(ctx, int(bobID),
		"Bob's Post", "Written by Bob", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Alice should see only her own post.
	alicePosts, err := svc.GetPostsByUser(ctx, int(aliceID))
	if err != nil {
		t.Fatal(err)
	}

	if len(alicePosts) != 1 {
		t.Fatalf("expected 1 Alice post, got %d", len(alicePosts))
	}

	if alicePosts[0].Title != "Alice's Post" {
		t.Fatalf("unexpected Alice post: %s", alicePosts[0].Title)
	}

	// Bob should see only his own post.
	bobPosts, err := svc.GetPostsByUser(ctx, int(bobID))
	if err != nil {
		t.Fatal(err)
	}

	if len(bobPosts) != 1 {
		t.Fatalf("expected 1 Bob post, got %d", len(bobPosts))
	}

	if bobPosts[0].Title != "Bob's Post" {
		t.Fatalf("unexpected Bob post: %s", bobPosts[0].Title)
	}

	t.Log("Alice sees only her own posts")
	t.Log("Bob sees only his own posts")
	t.Log("GetPostsByUser works correctly!")
}
func TestGetLikedPostsByUser(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.NewDataBase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.New(repo)

	// Create a user.
	result, err := db.ExecContext(ctx,
		`INSERT INTO users (username, email, password)
		 VALUES (?, ?, ?)`,
		"alice", "alice@example.com", "test-hash",
	)
	if err != nil {
		t.Fatal(err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create three posts.
	titles := []string{
		"Liked Post",
		"Disliked Post",
		"No Vote Post",
	}

	postIDs := make([]int64, 0, len(titles))

	for _, title := range titles {
		result, err := db.ExecContext(ctx,
			`INSERT INTO posts (user_id, title, body)
			 VALUES (?, ?, ?)`,
			userID, title, "Test content",
		)
		if err != nil {
			t.Fatal(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			t.Fatal(err)
		}

		postIDs = append(postIDs, id)
	}

	// Like the first post.
	_, err = db.ExecContext(ctx,
		`INSERT INTO post_like (user_id, post_id, vote)
		 VALUES (?, ?, 1)`,
		userID, postIDs[0],
	)
	if err != nil {
		t.Fatal(err)
	}

	// Dislike the second post.
	_, err = db.ExecContext(ctx,
		`INSERT INTO post_like (user_id, post_id, vote)
		 VALUES (?, ?, -1)`,
		userID, postIDs[1],
	)
	if err != nil {
		t.Fatal(err)
	}

	// Leave the third post without a vote.

	likedPosts, err := svc.GetLikedPostsByUser(ctx, int(userID))
	if err != nil {
		t.Fatal(err)
	}

	if len(likedPosts) != 1 {
		t.Fatalf("expected 1 liked post, got %d", len(likedPosts))
	}

	if likedPosts[0].ID != int(postIDs[0]) {
		t.Fatalf(
			"expected liked post ID %d, got %d",
			postIDs[0],
			likedPosts[0].ID,
		)
	}

	t.Log("Liked post returned correctly")
	t.Log("Disliked post excluded")
	t.Log("Post without vote excluded")
	t.Log("GetLikedPostsByUser works correctly!")
}
func TestSearchPosts(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.NewDataBase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.New(repo)

	// Create a test user.
	result, err := db.ExecContext(ctx,
		`INSERT INTO users (username, email, password)
		 VALUES (?, ?, ?)`,
		"alice", "alice@example.com", "test-hash",
	)
	if err != nil {
		t.Fatal(err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create three posts.
	testPosts := []struct {
		title   string
		content string
	}{
		{"The Dragon Kingdom", "A fantasy adventure"},
		{"Space Adventures", "Exploring distant galaxies"},
		{"My Favorite Books", "I really love dragons"},
	}

	for _, p := range testPosts {
		err := svc.CreatePost(ctx, int(userID), p.title, p.content, nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Search by title.
	posts, err := svc.SearchPosts(ctx, "Dragon Kingdom")
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 || posts[0].Title != "The Dragon Kingdom" {
		t.Fatalf("title search failed: got %+v", posts)
	}

	t.Log("Search by title works")

	// Search by body.
	posts, err = svc.SearchPosts(ctx, "galaxies")
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 || posts[0].Title != "Space Adventures" {
		t.Fatalf("body search failed: got %+v", posts)
	}

	t.Log("Search by body works")

	// Search with multiple matches.
	posts, err = svc.SearchPosts(ctx, "dragon")
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 dragon-related posts, got %d", len(posts))
	}

	t.Log("Search returns multiple matching posts")

	// Search with no matches.
	posts, err = svc.SearchPosts(ctx, "unicorn")
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 0 {
		t.Fatalf("expected 0 posts, got %d", len(posts))
	}

	t.Log("Unrelated posts are excluded")
	t.Log("SearchPosts works correctly!")
}
func TestDeletePost(t *testing.T) {
	ctx := context.Background()

	db, err := database.NewDataBase(
		filepath.Join(t.TempDir(), "delete_test.db"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.New(repo)

	// Create Alice.
	result, err := db.ExecContext(ctx,
		`INSERT INTO users (username, email, password)
         VALUES (?, ?, ?)`,
		"alice", "alice@example.com", "test-hash",
	)
	if err != nil {
		t.Fatal(err)
	}
	aliceID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create Bob.
	result, err = db.ExecContext(ctx,
		`INSERT INTO users (username, email, password)
         VALUES (?, ?, ?)`,
		"bob", "bob@example.com", "test-hash",
	)
	if err != nil {
		t.Fatal(err)
	}
	bobID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Alice creates a post.
	result, err = db.ExecContext(ctx,
		`INSERT INTO posts (user_id, title, body)
         VALUES (?, ?, ?)`,
		aliceID, "Alice's Post", "Test content",
	)
	if err != nil {
		t.Fatal(err)
	}
	postID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Bob cannot delete Alice's post", func(t *testing.T) {
		err := svc.DeletePost(ctx, postID, bobID)

		if !errors.Is(err, service.ErrUnauthorized) {
			t.Fatalf("expected unauthorized error, got %v", err)
		}

		// Verify the post still exists.
		_, err = repo.PostRepo.GetPostByID(ctx, postID)
		if err != nil {
			t.Fatalf("post should still exist: %v", err)
		}
	})

	t.Run("Alice can delete her own post", func(t *testing.T) {
		err := svc.DeletePost(ctx, postID, aliceID)
		if err != nil {
			t.Fatalf("expected successful deletion, got %v", err)
		}

		// Verify the post is gone.
		_, err = repo.PostRepo.GetPostByID(ctx, postID)

		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected post to be deleted, got %v", err)
		}
	})

	t.Run("deleting nonexistent post", func(t *testing.T) {
		err := svc.DeletePost(ctx, postID, aliceID)

		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows, got %v", err)
		}
	})
}
