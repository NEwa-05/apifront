package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var appPort = os.Getenv("APP_PORT")
var apiURL = os.Getenv("API_URL")

func getAPIHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Print(err)
	}
	apiIDValue := r.PostFormValue("id")
	JWTToken := r.PostFormValue("token")

	validURL, err := url.JoinPath(apiURL, apiIDValue)
	if err != nil {
		log.Print(err)
	}
	log.Printf("test: %s", validURL)
	apiResp, err := http.NewRequest("GET", validURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	if len(strings.TrimSpace(JWTToken)) != 0 {
		apiResp.Header.Add("Authorization", JWTToken)
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

func postAPIHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Print(err)
	}
	jsonBody := r.PostFormValue("data")
	JWTToken := r.PostFormValue("token")

	toto := bytes.NewBufferString(jsonBody)

	apiReq, err := http.NewRequest("POST", apiURL, toto)
	if err != nil {
		log.Fatalln(err)
	}

	if len(strings.TrimSpace(JWTToken)) != 0 {
		apiReq.Header.Add("Authorization", JWTToken)
	}

	client := &http.Client{}
	resp, err := client.Do(apiReq)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	a, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	tmpl := template.Must(template.New("tmpl").ParseFiles("postdata.html"))
	if err := tmpl.ExecuteTemplate(w, "postdata.html", string(a)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("tmpl").ParseFiles("index.html"))
	if err := tmpl.ExecuteTemplate(w, "index.html", "data"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/postdata", postAPIHandler)
	http.HandleFunc("/getdata", getAPIHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(appPort, nil)
}
