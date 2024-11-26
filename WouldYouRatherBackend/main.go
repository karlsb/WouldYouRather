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
var NUMBER_OF_PAIRS int = 0

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

//TODO: get number of pairs -> a function that returns the number of rows in the table pairs

func (db *Database) getNumberOfPairs() int {
	num := 0
	query := "SELECT COUNT(id) FROM pairs"
	err := db.sqldb.QueryRow(query).Scan(&num)
	if err != nil {
		log.Println("Error in executing getNumberOfRowsQuery: ", query)
	}

	return num
}

func createGetRandomPairQueryString(userID string) (string, []interface{}) {
	//	if len(SeenPairs[userID]) == NUMBER_OF_PAIRS {
	//SeenPairs = make(map[string][]int) //TODO bind this to a custom datastrucure for clarity
	////add allPairsSeen true to response
	//}
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
		log.Println("error in db.Exec: ", update_query, "| With ?  = ", choice.Id)
		log.Println("Db error massege:", err)
	}
	err = db.sqldb.QueryRow(select_query, choice.Id).Scan(&pair.Id, &pair.Left, &pair.Right, &pair.Lcount, &pair.Rcount)
	if err != nil {
		log.Println("error in db.QueryRow: ", select_query, " | With ? = ", choice.Id)
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
	Status       string   `json:"status"`
	Pair         TextPair `json:"pair"`
	AllPairsSeen bool     `json:"allPairsSeen"`
}

var db Database

func getRandomPairHandler(w http.ResponseWriter, r *http.Request) {
	//respond with json
	cookie, err := r.Cookie("user_id")
	if err != nil {
		newUserID := uuid.New().String()
		cookie = &http.Cookie{
			Name:     "user_id",
			Value:    newUserID,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Now().Add(24 * time.Hour),
		}
		log.Println("Generated new user ID:", newUserID)
		log.Println("cookie created: ", cookie)
	} else {
		log.Println("we have a cookie: ", cookie)
	}

	response := Message{
		Status:       "success",
		Pair:         db.getRandomPair(cookie.Value),
		AllPairsSeen: false,
	}

	if len(SeenPairs[cookie.Value]) == NUMBER_OF_PAIRS {
		SeenPairs = make(map[string][]int) //TODO bind this to a custom datastrucure for clarity
		//add allPairsSeen true to response
		response.AllPairsSeen = true
	}

	log.Println("getRandomPairHandler: ", SeenPairs)
	//TODO possibly move this to getRandonPair logic
	SeenPairs[cookie.Value] = append(SeenPairs[cookie.Value], response.Pair.Id)
	log.Println("getRandomPairHandler: ", SeenPairs)
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, cookie)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encdode JSON", http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the wouldyourather API")
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

	if cookie.Value == "" {
		log.Println("Cookie value is empty")
		http.Error(w, "Invalid cookie value", http.StatusBadRequest)
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

	log.Println("storeAnswerHandler - cookie value: ", cookie.Value)
	log.Println("storeAnsweHandler - map:", SeenPairs)
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

func handleGetNumberOfPairs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	type Integer struct {
		Value int `json:"value"`
	}
	var response Integer
	response.Value = db.getNumberOfPairs()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error in encoding response in handleGetNumberOfPairs: ", err)
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
			w.WriteHeader(http.StatusOK)
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
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	db.init()
	defer db.Close()
	NUMBER_OF_PAIRS = db.getNumberOfPairs()

	if NUMBER_OF_PAIRS == 0 {
		panic("NUMBER_OF_PAIRS SHOULD NOT BE 0. CHECK YOUR DATABASE")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", routeCheckMiddleware("/", indexHandler))
	mux.HandleFunc("/random-pair", routeCheckMiddleware("/random-pair", getRandomPairHandler))
	mux.HandleFunc("/store-answer", routeCheckMiddleware("/store-answer", storeAnswer))
	mux.HandleFunc("/get-number-of-pairs", routeCheckMiddleware("/get-number-of-pairs", handleGetNumberOfPairs))

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
