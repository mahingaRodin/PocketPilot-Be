package database

import (
    "database/sql"
    // "fmt"
    "log"
    "time"

  	 _"github.com/lib/pq"
)
	
type DB struct {
    *sql.DB
}

func Connect(databaseURL string) *DB {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    log.Println("Successfully connected to database!")
    return &DB{db}
}

func (db *DB) Close() error {
    return db.DB.Close()
}