package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Data struct {
	Quote string
}

var apiURL = os.Getenv("API_URL")

func apiHandler(w http.ResponseWriter, r *http.Request) {

	var result interface{}
	fmt.Printf("%s", r.Form["id"])
	apiIDValue := strings.Join(r.Form["id"], "")
	fmt.Println(apiIDValue)
	apiResp, err := http.Get(path.Join(apiURL, apiIDValue))
	fmt.Println(apiResp)
	if err != nil {
		log.Print(err)
	}
	err = json.NewDecoder(apiResp.Body).Decode(&result)
	if err != nil {
		log.Print(err)
	}

	a, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(a)
}

func main() {
	tmpl := template.Must(template.New("tmpl").ParseFiles("test.html"))
	appPort := os.Getenv("APP_PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "test.html", apiURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/ajax", apiHandler)
	http.ListenAndServe(appPort, nil)
}
