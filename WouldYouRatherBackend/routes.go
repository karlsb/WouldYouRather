package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

/*
*  API ROUTES
*
*
*
 */

type Response struct {
	Status       string   `json:"status"`
	Pair         TextPair `json:"pair"`
	AllPairsSeen bool     `json:"allPairsSeen"`
}

type MultiPairResponse struct {
	Status       string     `json:"status"`
	Pairs        []TextPair `json:"pair"`
	AllPairsSeen bool       `json:"allPairsSeen"`
}

var db Database

func handleGetNumberOfPairsN(w http.ResponseWriter, r *http.Request) {
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
	}

	response := MultiPairResponse{
		Status:       "success",
		Pairs:        db.getRandomPairN(cookie.Value, 4),
		AllPairsSeen: false,
	}

	if len(SeenPairs[cookie.Value]) == NUMBER_OF_PAIRS {
		SeenPairs = make(map[string][]int) //TODO bind this to a custom datastrucure for clarity
		response.AllPairsSeen = true
	}

	for _, pair := range response.Pairs {
		SeenPairs[cookie.Value] = append(SeenPairs[cookie.Value], pair.Id)
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, cookie)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encdode JSON", http.StatusInternalServerError)
	}

}

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
	}

	response := Response{
		Status:       "success",
		Pair:         db.getRandomPair(cookie.Value),
		AllPairsSeen: false,
	}

	if len(SeenPairs[cookie.Value]) == NUMBER_OF_PAIRS {
		SeenPairs = make(map[string][]int) //TODO bind this to a custom datastrucure for clarity
		response.AllPairsSeen = true
	}

	SeenPairs[cookie.Value] = append(SeenPairs[cookie.Value], response.Pair.Id)
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
		log.Println("ClientError in StoreAnswer: Cookie from client has no field user_id")
		http.Error(w, "Invalid cookie", http.StatusBadRequest)
	}

	if cookie.Value == "" {
		log.Println("ClientError in StoreAnswer: Cookie value from client is empty")
		http.Error(w, "Invalid cookie", http.StatusBadRequest)
		return
	}

	var choice Choice

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&choice); err != nil {
		log.Println("ClientError in StoreAnswer: Failed to Decode JSON")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if choice.LeftRight != "right" && choice.LeftRight != "left" {
		log.Println("ClientError in StoreAnswer: Invalid value in Choice field in client JSON")
		http.Error(w, "Invalid value of field leftright", http.StatusBadRequest)
		return
	}

	response := Response{
		Status: "success",
		Pair:   db.increaseCountAndReturnPair(choice),
	}

	_, exists := SeenPairs[cookie.Value]
	if !exists {
		log.Println("ClientError in StoreAnswer: Invalid user id in cookie.Value. Client provided: ", cookie.Value)
		http.Error(w, "Invalid Cookie", http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("ServerError in StoreAnswer: Failed to encode JSON", cookie.Value)
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
		log.Println("ServerError in encoding response in handleGetNumberOfPairs: ", err)
	}
}

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
