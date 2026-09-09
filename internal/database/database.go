package database

import "database/sql"

const schema = `  CREATE TABLE IF NOT EXISTS user (
				id 	INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT NOT NULL UNIQUE,
    			email TEXT NOT NULL UNIQUE,
				password TEXT nOT NULL,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS sessions (
				uuid INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				expires_at DATETIME NOT NULL,

				FOREIGN KEY (user_id) REFERENCES users(id)
				);

				CREATE TABLE IF NOT EXISTS posts (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				title TEXT NOT NULL,
				body TEXT NOT NULL,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

				FOREIGN KEY (user_id) REFERENCES users(id)
				);

				CREATE TABLE IF NOT EXISTS comments (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			post_id INTEGER NOT NULL,
    			user_id INTEGER NOT NULL,
    			content TEXT NOT NULL,
    			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    			FOREIGN KEY (post_id) REFERENCES posts(id),
    			FOREIGN KEY (user_id) REFERENCES users(id)
				);

				CREATE TABLE categories (
   				id INTEGER PRIMARY KEY AUTOINCREMENT,
    			name TEXT NOT NULL UNIQUE
				);


				CREATE TABLE post_like (
  			    user_id INTEGER NOT NULL,
    			post_id INTEGER NOT NULL,
    			like INTEGER NOT NULL,

    			PRIMARY KEY (user_id, post_id),

    			CHECK (vote IN (-1, 1)),

    			FOREIGN KEY (user_id) REFERENCES users(id),
    			FOREIGN KEY (post_id) REFERENCES posts(id)
				);

				CREATE TABLE comment_like (
  			    user_id INTEGER NOT NULL,
				comment_id INTEFER NOT NULL
    			like INTEGER NOT NULL,

    			PRIMARY KEY (user_id, post_id, comment_id),

    			CHECK (vote IN (-1, 1)),

    			FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (comment_id) REFERENCES comments(id)
				);


				`

func NewDataBase(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// InitializeQuery creates the database tables using the defined schema.
func InitializeQuery(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
