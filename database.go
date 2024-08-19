package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tursodatabase/go-libsql"
)

func getChannels() {
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
	rows, err := db.Query("SELECT * FROM channels;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	// var channels []Channel

	for rows.Next() {
		var channel Channel

		if err := rows.Scan(&channel.ID, &channel.Title); err != nil {
			fmt.Println("Error scanning row:", err)
		}

		// channels = append(channels, channel)
		fmt.Println(channel.ID, channel.Title)
	}

	// fmt.Println(len(channels))

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
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

	return "Registered successfully!"
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
