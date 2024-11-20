package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // fast but be careful, since it is written in CGO - might have compatibility issues
)

func main() {
	fmt.Println("hello from builddbscript")
	db, err := sql.Open("sqlite", "./wouldyourather.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS pairs (
            id INTEGER PRIMARY KEY,
            left TEXT NOT NULL,
			right TEXT NOT NULL,
			lcount INTEGER DEFAULT 0 NOT NULL,
			rcount INTEGER DEFAULT 0 NOT NULL
        );
    `)
	if err != nil {
		fmt.Println("failed to run migrations: %w", err)
	}

	// Seed data
	_, err = db.Exec(`INSERT INTO pairs (left,right) VALUES ('Would you', 'rather')`)
	if err != nil {
		fmt.Println("Data already exists or failed to seed")
	} else {
		fmt.Println("Database seeded successfully")
	}
}
