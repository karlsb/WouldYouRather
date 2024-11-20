package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite" // fast but be careful, since it is written in CGO - might have compatibility issues
)

type TextPair struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

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

func (db Database) Close() {
	db.sqldb.Close()
}

// NOW ONLY GETS TOP
func (db Database) getRandomPair() TextPair {
	// get number of rows pick an id in the range of num of rows
	rows, err := db.sqldb.Query("SELECT left, right FROM pairs ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		fmt.Println("error in db.Query")
	}
	defer rows.Close()

	var left string
	var right string
	for rows.Next() {
		rows.Scan(&left, &right)
	}
	return TextPair{Left: left, Right: right}
}

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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Origin", "https://whatwouldyourather.netlify.app") //TODO CHANGE HERE IF PRODUCTION
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader((http.StatusNoContent))
			return
		}
		next.ServeHTTP(w, r)
	})

}

func main() {
	db.init()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", routeCheckMiddleware("/", indexHandler))
	mux.HandleFunc("/random-pair", routeCheckMiddleware("/random-pair", getRandomPairHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, enableCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
