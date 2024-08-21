package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tursodatabase/go-libsql"
)

func findVideoId(substr string) (string, error) {
	// Set up the database
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	dbName := "local.db"
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	syncInterval := time.Minute

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)

	if err != nil {
		return "", err
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	// Do something with the database
	rows, err := db.Query("SELECT video_id FROM posts where video_id like ?;", "%"+substr)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post

		if err := rows.Scan(&post.video_id); err != nil {
			return "", err
		}

		return post.video_id, nil
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func hasRegistered(id int64) bool {
	// Set up the database
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	dbName := "local.db"
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	syncInterval := time.Minute

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)

	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	// Do something with the database
	rows, err := db.Query("SELECT * FROM channels WHERE id = ?;", id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	return rows.Next()
}

func registerChannel(id int64, title string, password string) string {
	// Set up the database
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	dbName := "local.db"
	dir, err := os.MkdirTemp("", "libsql-*")

	if err != nil {
		return "Error creating temporary directory!"
	}

	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	syncInterval := time.Minute

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)

	if err != nil {
		return "Error creating db connector!"
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	hash, err := HashPassword(password)
	if err != nil {
		return "Error hashing password!"
	}

	// Do something with the database
	_, err = db.Query("INSERT INTO channels (id, title, password) VALUES (?, ? , ?)", id, title, hash)
	if err != nil {
		return "Error registering channel!"
	}

	return "Registered successfully! You can now log in."
}

func loginChannel(id int64, password string) string {
	// Set up the database
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	dbName := "local.db"
	dir, err := os.MkdirTemp("", "libsql-*")

	if err != nil {
		return "Error creating temporary directory!"
	}

	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	syncInterval := time.Minute

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)

	if err != nil {
		return "Error creating db connector!"
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	// Do something with the database
	rows, err := db.Query("Select * FROM channels WHERE id = ?;", id)
	if err != nil {
		return "Error logging in channel!"
	}
	defer rows.Close()

	var channel Channel
	if rows.Next() {
		if err := rows.Scan(&channel.ID, &channel.Title, &channel.Password); err != nil {
			return "Error scanning row"
		}
	}
	if err := rows.Err(); err != nil {
		return "Error during rows iteration"
	}

	if !VerifyPassword(password, channel.Password) {
		return "Wrong password!"
	}

	return ""
}

func createPost(id int64, title string, rating int, description string, tags string, video_id string, video_path string, thumbnail_id string, thumbnail_path string, g_thumbnail_id string, g_thumbnail_path string) string {
	// Set up the database
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	dbName := "local.db"
	dir, err := os.MkdirTemp("", "libsql-*")

	if err != nil {
		return "Error creating temporary directory!"
	}

	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	syncInterval := time.Minute

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)

	if err != nil {
		return "Error creating db connector!"
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	// Do something with the database
	_, err = db.Query("INSERT INTO posts VALUES (?, ? , ?, ?, ?, ?, ?, ?, ?, ?, ?)", id, title, rating, description, tags, video_id, video_path, thumbnail_id, thumbnail_path, g_thumbnail_id, g_thumbnail_path)
	if err != nil {
		return "Error creating post!"
	}

	return "Post created successfully!"
}