package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
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

// -----UUID: string - pairIDs: []int
var SeenPairs = make(map[string][]int)

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

func createGetRandomPairQueryString(userID string) (string, []interface{}) {
	placeholders := make([]string, len(SeenPairs[userID]))
	args := make([]interface{}, len(SeenPairs[userID]))
	for i, id := range SeenPairs[userID] {
		placeholders[i] = "?"
		args[i] = id
	}

	// Create the query string
	queryString := fmt.Sprintf(
		"SELECT id, left, right, lcount, rcount FROM pairs WHERE id NOT IN (%s) ORDER BY RANDOM() LIMIT 1",
		strings.Join(placeholders, ","),
	)
	return queryString, args
}

// NOW ONLY GETS TOP
func (db Database) getRandomPair(userID string) TextPair {
	// get number of rows pick an id in the range of num of rows
	queryString, args := createGetRandomPairQueryString(userID)
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
	cookie, err := r.Cookie("user_id")
	if err != nil {
		newUserID := uuid.New().String()
		cookie = &http.Cookie{
			Name:    "user_id",
			Value:   newUserID,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		log.Println("Generated new user ID:", newUserID)
		log.Println("cookie created: ", cookie)
	} else {
		log.Println("we have a cookie: ", cookie)
	}

	response := Message{
		Status: "success",
		Pair:   db.getRandomPair(cookie.Value),
	}

	SeenPairs[cookie.Value] = append(SeenPairs[cookie.Value], response.Pair.Id)
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, cookie)
	log.Println("Headers", w.Header())
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encdode JSON", http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello There")
}

func storeAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Invalid cookie", http.StatusBadRequest)
	}

	log.Println("cookie: ", cookie)

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

	value, exists := SeenPairs[cookie.Value]
	if !exists {
		log.Fatal("we dont have a SeenPair value of: ", cookie.Value)
	} else {
		log.Println(value)
	}

	// TODO I can break this out into its own method
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encdode JSON", http.StatusInternalServerError)
	}

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
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader((http.StatusNoContent))
			return
		}
		next.ServeHTTP(w, r)
	})

}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log_level := os.Getenv("LOG_LEVEL")
		//if log_level == "DEV" {
		//log.Println(r.Body)
		//}
		next.ServeHTTP(w, r)
	})

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.init()
	defer db.Close()

	//TESTTING

	//

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
