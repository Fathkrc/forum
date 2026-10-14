package handler

import (
	"cmd/forum/internal/database"
	"cmd/forum/internal/repository"
	"cmd/forum/internal/service"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearchPostsHandler(t *testing.T) {
	// 1. Create a temporary SQLite database.
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.NewDataBase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// 2. Initialize the real repository and service.
	repo := repository.NewRepository(db)
	svc := service.New(repo)

	// 3. Create our handler and register the real routes.
	h := New(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	ctx := context.Background()

	// 4. Create a test user.
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

	// 5. Create test posts.
	testPosts := []struct {
		title   string
		content string
	}{
		{"The Dragon Kingdom", "A fantasy adventure"},
		{"Space Adventures", "Exploring distant galaxies"},
		{"My Favorite Books", "I really love dragons"},
	}

	for _, p := range testPosts {
		err := svc.CreatePost(
			ctx,
			int(userID),
			p.title,
			p.content,
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	// 6. Define our HTTP test cases.
	tests := []struct {
		name           string
		url            string
		wantStatus     int
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:       "search matching posts",
			url:        "/posts/search?q=dragon",
			wantStatus: http.StatusOK,
			wantContains: []string{
				"The Dragon Kingdom",
				"My Favorite Books",
			},
			wantNotContain: []string{
				"Space Adventures",
			},
		},
		{
			name:         "search with no matches",
			url:          "/posts/search?q=unicorn",
			wantStatus:   http.StatusOK,
			wantContains: []string{"No matching posts found."},
		},
		{
			name:         "empty search query",
			url:          "/posts/search?q=",
			wantStatus:   http.StatusBadRequest,
			wantContains: []string{"Search query is required"},
		},
	}

	// 7. Run each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d. Body: %s",
					tt.wantStatus,
					rec.Code,
					rec.Body.String(),
				)
			}

			body := rec.Body.String()

			for _, expected := range tt.wantContains {
				if !strings.Contains(body, expected) {
					t.Errorf(
						"expected response to contain %q, got %q",
						expected,
						body,
					)
				}
			}

			for _, unexpected := range tt.wantNotContain {
				if strings.Contains(body, unexpected) {
					t.Errorf(
						"response should not contain %q, got %q",
						unexpected,
						body,
					)
				}
			}
		})
	}
}
func TestListPostsByCategoryHandler(t *testing.T) {
	ctx := context.Background()

	// Temporary database
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.NewDataBase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.New(repo)

	h := New(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Create a test user
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

	// Create two categories
	result, err = db.ExecContext(ctx,
		`INSERT INTO categories (name) VALUES (?)`,
		"Fantasy",
	)
	if err != nil {
		t.Fatal(err)
	}

	fantasyID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	result, err = db.ExecContext(ctx,
		`INSERT INTO categories (name) VALUES (?)`,
		"Science Fiction",
	)
	if err != nil {
		t.Fatal(err)
	}

	scifiID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create posts in different categories
	err = svc.CreatePost(
		ctx,
		int(userID),
		"The Dragon Castle",
		"A fantasy adventure",
		[]int{int(fantasyID)},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = svc.CreatePost(
		ctx,
		int(userID),
		"Journey to Mars",
		"A science fiction adventure",
		[]int{int(scifiID)},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Test category filtering through HTTP
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/posts/category?category=%d", fantasyID),
		nil,
	)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	body := rec.Body.String()

	if !strings.Contains(body, "The Dragon Castle") {
		t.Error("expected Fantasy post in response")
	}

	if strings.Contains(body, "Journey to Mars") {
		t.Error("Science Fiction post should not appear")
	}
}
func setupPostTest(t *testing.T) (*sql.DB, *service.Service, *http.ServeMux) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.NewDataBase(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repo := repository.NewRepository(db)
	svc := service.New(repo)

	h := New(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	return db, svc, mux
}
func TestListPostsByUserHandler(t *testing.T) {
	db, svc, mux := setupPostTest(t)
	ctx := context.Background()

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

	// Create one post for each user.
	if err := svc.CreatePost(ctx, int(aliceID),
		"Alice's Fantasy Post", "Written by Alice", nil); err != nil {
		t.Fatal(err)
	}

	if err := svc.CreatePost(ctx, int(bobID),
		"Bob's Space Post", "Written by Bob", nil); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		url            string
		wantStatus     int
		wantContains   string
		wantNotContain string
	}{
		{
			name:           "Alice's posts only",
			url:            fmt.Sprintf("/posts/user?user=%d", aliceID),
			wantStatus:     http.StatusOK,
			wantContains:   "Alice's Fantasy Post",
			wantNotContain: "Bob's Space Post",
		},
		{
			name:           "Bob's posts only",
			url:            fmt.Sprintf("/posts/user?user=%d", bobID),
			wantStatus:     http.StatusOK,
			wantContains:   "Bob's Space Post",
			wantNotContain: "Alice's Fantasy Post",
		},
		{
			name:         "invalid user ID",
			url:          "/posts/user?user=abc",
			wantStatus:   http.StatusBadRequest,
			wantContains: "Invalid user ID",
		},
		{
			name:         "missing user ID",
			url:          "/posts/user",
			wantStatus:   http.StatusBadRequest,
			wantContains: "Invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d: %s",
					tt.wantStatus, rec.Code, rec.Body.String())
			}

			body := rec.Body.String()

			if !strings.Contains(body, tt.wantContains) {
				t.Errorf("expected response to contain %q, got %q",
					tt.wantContains, body)
			}

			if tt.wantNotContain != "" && strings.Contains(body, tt.wantNotContain) {
				t.Errorf("response should not contain %q", tt.wantNotContain)
			}
		})
	}
}
func TestListLikedPostsHandler(t *testing.T) {
	db, _, mux := setupPostTest(t)
	ctx := context.Background()

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
	titles := []string{
		"Liked Dragon Post",
		"Disliked Space Post",
		"Unvoted Fantasy Post",
	}

	var postIDs []int64

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

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "liked posts",
			url:        fmt.Sprintf("/posts/liked?user=%d", userID),
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid user ID",
			url:        "/posts/liked?user=abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing user ID",
			url:        "/posts/liked",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d: %s",
					tt.wantStatus,
					rec.Code,
					rec.Body.String(),
				)
			}

			if tt.name == "liked posts" {
				body := rec.Body.String()

				if !strings.Contains(body, "Liked Dragon Post") {
					t.Error("liked post should appear")
				}

				if strings.Contains(body, "Disliked Space Post") {
					t.Error("disliked post should not appear")
				}

				if strings.Contains(body, "Unvoted Fantasy Post") {
					t.Error("unvoted post should not appear")
				}
			}
		})
	}
}
