package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var appPort = os.Getenv("APP_PORT")
var apiURL = os.Getenv("API_URL")

func apiHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Print(err)
	}
	apiIDValue := r.PostFormValue("id")

	validURL, err := url.JoinPath(apiURL, apiIDValue)
	if err != nil {
		log.Print(err)
	}

	apiResp, err := http.NewRequest("GET", validURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	apiResp.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(apiResp)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	a, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	tmpl := template.Must(template.New("tmpl").ParseFiles("data.html"))
	if err := tmpl.ExecuteTemplate(w, "data.html", string(a)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("tmpl").ParseFiles("index.html"))
	tmpl.Delims("[[", "]]")
	if err := tmpl.ExecuteTemplate(w, "index.html", "data"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/data", apiHandler)
	http.ListenAndServe(appPort, nil)
}
