package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

/*
*  DATABASE
*
*
*
 */
type Database struct {
	sqldb *sql.DB
}

func (db *Database) init() {
	temp, err := sql.Open("sqlite", "./build-database/wouldyourather.db")
	if err != nil {
		log.Fatal("Failed to Open Database", err)
	}
	db.sqldb = temp
}

func (db *Database) Close() {
	db.sqldb.Close()
}

func (db *Database) getNumberOfPairs() int {
	num := 0
	query := "SELECT COUNT(id) FROM pairs"
	err := db.sqldb.QueryRow(query).Scan(&num)
	if err != nil {
		log.Println("Error in executing getNumberOfRowsQuery: ", query)
	}

	return num
}

func createQueryStringGetRandomPair(userID string) (string, []interface{}) {
	placeholders := make([]string, len(SeenPairs[userID]))
	args := make([]interface{}, len(SeenPairs[userID]))
	for i, id := range SeenPairs[userID] {
		placeholders[i] = "?"
		args[i] = id
	}

	queryString := fmt.Sprintf(
		"SELECT id, left, right, lcount, rcount FROM pairs WHERE id NOT IN (%s) ORDER BY RANDOM() LIMIT 1",
		strings.Join(placeholders, ","),
	)
	return queryString, args
}

func (db Database) getRandomPair(userID string) TextPair {
	queryString, args := createQueryStringGetRandomPair(userID)
	rows, err := db.sqldb.Query(queryString, args...)
	if err != nil {
		fmt.Println("error in db.Query")
	}
	defer rows.Close()

	var pair TextPair

	for rows.Next() {
		rows.Scan(&pair.Id, &pair.Left, &pair.Right, &pair.Lcount, &pair.Rcount)
	}
	return pair
}

func (db Database) increaseCountAndReturnPair(choice Choice) TextPair {
	update_query := "UPDATE pairs SET lcount = lcount + 1 WHERE id = ?"
	if choice.LeftRight == "right" {
		update_query = "UPDATE pairs SET rcount = rcount + 1 WHERE id = ?"
	}

	select_query := "SELECT * FROM pairs WHERE id = ? LIMIT 1"

	var pair TextPair

	_, err := db.sqldb.Exec(update_query, choice.Id)
	if err != nil {
		log.Println("error in db.Exec: ", update_query, "| With ?  = ", choice.Id, " - ", err)
	}
	err = db.sqldb.QueryRow(select_query, choice.Id).Scan(&pair.Id, &pair.Left, &pair.Right, &pair.Lcount, &pair.Rcount)
	if err != nil {
		log.Println("Error in db.QueryRow: ", select_query, " | With ? = ", choice.Id)
	}

	return pair
}
