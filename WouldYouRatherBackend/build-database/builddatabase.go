package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite" // fast but be careful, since it is written in CGO - might have compatibility issues
)

func connectToDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createPairsTable(db *sql.DB) {
	// Create table
	_, err := db.Exec(`
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

}

func populateDBfromCSV(db *sql.DB, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	transaction, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	stmt, err := transaction.Prepare("INSERT INTO pairs (left, right) VALUES (?,?)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, row := range rows[1:] { //skip first line
		_, err = stmt.Exec(row[0], row[1])
		if err != nil {
			transaction.Rollback()
			log.Fatalf("Failed to execute statement: %v", err)
		}
	}

	if err := transaction.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Database populated successfully.")
}

func main() {

	filename := "data.csv"

	db := connectToDB("./wouldyourather.db")
	defer db.Close()

	createPairsTable(db)

	populateDBfromCSV(db, filename)

	// Seed data
}
