package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var appPort = os.Getenv("APP_PORT")

var Lines []Line

type Line struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	StartStation   string `json:"startStation,omitempty"`
	EndStation     string `json:"endStation,omitempty"`
	NumberStations int32  `json:"numberStations,omitempty"`
}

func postAPIHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var line Line
	json.Unmarshal(reqBody, &line)
	Lines = append(Lines, line)
	json.NewEncoder(w).Encode(line)

}

func returnSingleLine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, line := range Lines {
		if line.ID == key {
			json.NewEncoder(w).Encode(line)
		}
	}

}

func returnAllLines(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllEpisodes")
	json.NewEncoder(w).Encode(Lines)
}

func main() {
	Lines = []Line{
		Line{ID: "1", Name: "Ligne 1", StartStation: "La Défense - Grande Arche", EndStation: "Château de Vincennes", NumberStations: 25},
		Line{ID: "2", Name: "Ligne 2", StartStation: "Porte Dauphine", EndStation: "Nation", NumberStations: 25},
	}
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/newlines", returnAllLines)
	myRouter.HandleFunc("/newline", postAPIHandler).Methods("POST")
	myRouter.HandleFunc("/newline/{id}", returnSingleLine)
	log.Fatal(http.ListenAndServe(appPort, myRouter))
}
