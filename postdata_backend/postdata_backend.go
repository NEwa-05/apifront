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

var Episodes []Episode

type Episode struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	SerialID          int32   `json:"serialID"`
	EpisodeOrder      string  `json:"episodeOrder,omitempty"`
	OriginalAirDate   string  `json:"originalAirDate,omitempty"`
	Runtime           string  `json:"runtime,omitempty"`
	UKViewersMM       float32 `json:"ukViewersMM,omitempty"`
	AppreciationIndex float32 `json:"appreciationIndex,omitempty"`
}

func postAPIHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var episode Episode
	json.Unmarshal(reqBody, &episode)
	Episodes = append(Episodes, episode)
	json.NewEncoder(w).Encode(episode)

}

func returnSinglepisode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, episode := range Episodes {
		if episode.ID == key {
			json.NewEncoder(w).Encode(episode)
		}
	}

}

func returnAllEpisodes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllEpisodes")
	json.NewEncoder(w).Encode(Episodes)
}

func main() {
	Episodes = []Episode{}
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/newepisodes", returnAllEpisodes)
	myRouter.HandleFunc("/newepisode", postAPIHandler).Methods("POST")
	myRouter.HandleFunc("/newepisode/{id}", returnSinglepisode)
	log.Fatal(http.ListenAndServe(appPort, myRouter))
}
