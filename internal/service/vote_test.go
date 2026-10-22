package service_test

import (
	"cmd/forum/internal/database"
	"cmd/forum/internal/repository"
	"cmd/forum/internal/service"
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
)

func TestVotePost(t *testing.T) {
	ctx := context.Background()

	db, err := database.NewDataBase(
		filepath.Join(t.TempDir(), "vote_test.db"),
	)
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

	// Create a post.
	result, err = db.ExecContext(ctx,
		`INSERT INTO posts (user_id, title, body)
		 VALUES (?, ?, ?)`,
		userID, "Voting Test Post", "Test body",
	)
	if err != nil {
		t.Fatal(err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		isLike   bool
		wantVote int
	}{
		{"like post", true, 1},
		{"change to dislike", false, -1},
		{"change back to like", true, 1},
		{"repeat like", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.VotePost(
				ctx,
				int(userID),
				int(postID),
				tt.isLike,
			)
			if err != nil {
				t.Fatal(err)
			}

			got, err := repo.VoteRepo.GetPostVote(
				ctx,
				userID,
				postID,
			)

			// No vote should exist after clicking like twice.
			if tt.wantVote == 0 {
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf(
						"expected vote to be removed, got vote=%d, err=%v",
						got, err,
					)
				}
				return
			}

			// For other cases, a vote should exist.
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.wantVote {
				t.Errorf(
					"expected vote %d, got %d",
					tt.wantVote,
					got,
				)
			}
		})
	}

	// Ensure repeated votes didn't create duplicate rows.
	var count int

	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*)
		 FROM post_like
		 WHERE user_id = ? AND post_id = ?`,
		userID, postID,
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Errorf("expected 0 vote rows, got %d", count)
	}
}
