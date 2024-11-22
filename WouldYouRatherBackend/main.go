package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite" // fast but be careful, since it is written in CGO - might have compatibility issues
)

type TextPair struct {
	Id     int    `json:"id"`
	Left   string `json:"left"`
	Right  string `json:"right"`
	Lcount int    `json:"lcount"`
	Rcount int    `json:"rcount"`
}

type Choice struct {
	Id        int    `json:"id"`
	LeftRight string `json:"leftright"`
}

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
		log.Fatal(err)
	}
	db.sqldb = temp
}

func (db *Database) Close() {
	db.sqldb.Close()
}

// NOW ONLY GETS TOP
func (db Database) getRandomPair() TextPair {
	// get number of rows pick an id in the range of num of rows
	//New comment for testing
	rows, err := db.sqldb.Query("SELECT id, left, right, lcount, rcount FROM pairs ORDER BY RANDOM() LIMIT 1")
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
		fmt.Println("error in db.Exec: ", update_query, "| With ?  = ", choice.Id)
	}
	err = db.sqldb.QueryRow(select_query, choice.Id).Scan(&pair.Id, &pair.Left, &pair.Right, &pair.Lcount, &pair.Rcount)
	if err != nil {
		fmt.Println("error in db.QueryRow: ", select_query, " | With ? = ", choice.Id)
	}

	return pair
}

/*
*  API ROUTES
*
*
*
 */

type Message struct {
	Status string   `json:"status"`
	Pair   TextPair `json:"pair"`
}

var db Database

func getRandomPairHandler(w http.ResponseWriter, r *http.Request) {
	//respond with json

	response := Message{
		Status: "success",
		Pair:   db.getRandomPair(),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encdode JSON", http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

// a template for middleware, can use for logging aswell?
func routeCheckMiddleware(expectedPath string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedPath {
			http.NotFound(w, r)
			return
		}
		handler(w, r)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed_origin := os.Getenv("ALLOWED_ORIGIN")
		w.Header().Set("Access-Control-Allow-Origin", allowed_origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader((http.StatusNoContent))
			return
		}
		next.ServeHTTP(w, r)
	})

}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log_level := os.Getenv("LOG_LEVEL")
		if log_level == "DEV" {
			log.Println(r.Body)
		}
		next.ServeHTTP(w, r)
	})

}

func storeAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var choice Choice

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&choice); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if choice.LeftRight != "right" && choice.LeftRight != "left" {
		http.Error(w, "Invalid value of field leftright", http.StatusBadRequest)
		return
	}

	response := Message{
		Status: "success",
		Pair:   db.increaseCountAndReturnPair(choice),
	}

	// TODO I can break this out into its own method
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encdode JSON", http.StatusInternalServerError)
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.init()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", routeCheckMiddleware("/", indexHandler))
	mux.HandleFunc("/random-pair", routeCheckMiddleware("/random-pair", getRandomPairHandler))
	mux.HandleFunc("/store-answer", routeCheckMiddleware("/store-answer", storeAnswer))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, loggerMiddleware(enableCORS(mux))); err != nil {
		log.Fatal(err)
	}
}
