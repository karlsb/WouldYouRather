package main

import (
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

// -----UUID: string - pairIDs: []int
var SeenPairs = make(map[string][]int) //TODO should i not just attach this value to Database struct?
var NUMBER_OF_PAIRS int = 0            //TODO should i not just attach this value to Database struct?

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
		log.Println("WARNING: NUMBER_OF_PAIRS SHOULD NOT BE 0. CHECK YOUR DATABASE")
	}

	mux := http.NewServeMux()

	//TODO Change route names store-answer -> choice/answer
	//remove get-number-of-pairs or change to number-of-pairs
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
	if err := http.ListenAndServe(":"+port, enableCORS(mux)); err != nil {
		log.Fatal("ERROR in ListenAndServe, check your provided port number", err)
	}
}
