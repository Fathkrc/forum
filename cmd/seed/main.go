package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type book struct{ title, category string }

func main() {
	db, err := sql.Open("sqlite3", "data/forum.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run `make run` once first to create the database and categories.
	var categoryCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&categoryCount); err != nil || categoryCount == 0 {
		log.Fatal("Database not initialized. Run make run once, then stop it before seeding.")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	names := []string{"lion_reader_demo", "chapter_hunter_demo", "page_turner_demo"}
	userIDs := make([]int64, len(names))
	for i, name := range names {
		_, err = tx.Exec("INSERT OR IGNORE INTO users (username, email, password) VALUES (?, ?, ?)", name, name+"@example.invalid", "DEMO_ONLY_NO_LOGIN")
		if err != nil {
			log.Fatal(err)
		}
		if err = tx.QueryRow("SELECT id FROM users WHERE username = ?", name).Scan(&userIDs[i]); err != nil {
			log.Fatal(err)
		}
	}

	books := []book{
		{"Dune", "Science Fiction"}, {"The Hobbit", "Fantasy"}, {"1984", "Literary Fiction"},
		{"Pride and Prejudice", "Romance"}, {"The Great Gatsby", "Classics"},
		{"The Hunger Games", "Adventure"}, {"Frankenstein", "Horror"},
		{"The Lord of the Rings", "Fantasy"}, {"Jane Eyre", "Classics"},
		{"The Martian", "Science Fiction"}, {"Little Women", "Classics"},
		{"The Book Thief", "Historical Fiction"}, {"Dracula", "Horror"},
		{"The Name of the Wind", "Fantasy"}, {"To Kill a Mockingbird", "Classics"},
		{"The Handmaid's Tale", "Literary Fiction"}, {"The Alchemist", "Literary Fiction"},
		{"The Midnight Library", "Literary Fiction"}, {"The Silent Patient", "Mystery"},
		{"Project Hail Mary", "Science Fiction"},
	}
	titles := []string{"What did you think of %s?", "Favorite character in %s?", "Is %s worth reading?", "Best moment in %s?", "Books similar to %s?"}
	bodies := []string{
		"I just finished %s. What did you like most, and what would you change?",
		"Which character from %s stayed with you after reading? Why?",
		"I'm thinking about reading %s. Would you recommend it without spoilers?",
		"Which scene from %s made the biggest impression on you? Please mark spoilers.",
		"I enjoyed %s and want something similar. Any recommendations?",
	}
	comments := []string{
		"I loved the atmosphere and characters!", "I had a different take on the ending.",
		"This is going straight onto my reading list.", "The writing style was my favorite part.",
		"I would definitely recommend it.", "I think the middle chapters were the strongest.",
	}

	for i := 0; i < 100; i++ {
		b := books[i%len(books)]
		title := fmt.Sprintf(titles[i/len(books)], b.title)
		body := fmt.Sprintf(bodies[i/len(books)], b.title)
		author := userIDs[i%len(userIDs)]

		var postID int64
		err = tx.QueryRow("SELECT id FROM posts WHERE user_id = ? AND title = ? AND body = ? LIMIT 1", author, title, body).Scan(&postID)
		if err == sql.ErrNoRows {
			result, e := tx.Exec("INSERT INTO posts (user_id, title, body) VALUES (?, ?, ?)", author, title, body)
			if e != nil {
				log.Fatal(e)
			}
			postID, err = result.LastInsertId()
		}
		if err != nil {
			log.Fatal(err)
		}

		var categoryID int64
		err = tx.QueryRow("SELECT id FROM categories WHERE lower(name) = ?", strings.ToLower(b.category)).Scan(&categoryID)
		if err != nil {
			log.Fatalf("Missing category %q: %v", b.category, err)
		}
		if _, err = tx.Exec("INSERT OR IGNORE INTO post_category (post_id, category_id) VALUES (?, ?)", postID, categoryID); err != nil {
			log.Fatal(err)
		}

		for j := 0; j < 2; j++ {
			commenter := userIDs[(i+j+1)%len(userIDs)]
			content := comments[(i+j)%len(comments)]
			var commentID int64
			err = tx.QueryRow("SELECT id FROM comments WHERE post_id = ? AND user_id = ? AND content = ? LIMIT 1", postID, commenter, content).Scan(&commentID)
			if err == sql.ErrNoRows {
				result, e := tx.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, commenter, content)
				if e != nil {
					log.Fatal(e)
				}
				commentID, err = result.LastInsertId()
			}
			if err != nil {
				log.Fatal(err)
			}
			if _, err = tx.Exec("INSERT OR IGNORE INTO comment_like (user_id, comment_id, vote) VALUES (?, ?, 1)", author, commentID); err != nil {
				log.Fatal(err)
			}
		}

		for j := 1; j <= 2; j++ {
			voter := userIDs[(i+j)%len(userIDs)]
			vote := 1
			if j == 2 && i%4 == 0 {
				vote = -1
			}
			if _, err = tx.Exec("INSERT OR IGNORE INTO post_like (user_id, post_id, vote) VALUES (?, ?, ?)", voter, postID, vote); err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Seed complete: 3 demo users, 100 posts, 200 comments, and votes. Safe to rerun.")
}
