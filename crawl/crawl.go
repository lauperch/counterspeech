package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Text struct {
	Content string
	URL     string
	Source  string
	IsHS    bool
	IsNotHS bool
	Idk     bool
}

type Source struct {
	domain    string
	startUrl  string
	isRunning bool
}

// TODO global map is not v nice ;)
var sources = map[string]Source{}

const allowedDomains = [1]string{"20min.ch"}

func Status(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	domain := ps.ByName("domain")
	responseJSON(w, "")
}

func Run(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	domain := ps.ByName("runDomain")
	if contains(allowedDomains, domain) {
		responseJSON(w, "domain not yet implemented")
	}
	if sources[domain] != nil && sources["domain"].isRunning {
		responseJSON(w, "source already running")
	}
	src := Source{domain: domain}
	go Scrape(src)
	responseJSON(w, "")
}

func Stop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	domain := ps.ByName("stopDomain")
	responseJSON(w, "")
}

func Scrape(src Source) {

}

func NextUrl(prevUrl string) string {
	// TODO implement
	return ""
}

func HtmlToText(html string) Text {
	// TODO implement
	return nil
}

func Save(text Text) {
	url := ""
	if os.Getenv("APP_ENV") == "prod" {
		url := "" // TODO insert correct prod url
	} else {
		url := "http://localhost:5000"
	}
	resp, err := http.Post(url, "application/json", json.Marshal(text))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func MakeRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func setCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func contains(slice []string, element string) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

func main() {
	router := httprouter.New()
	router.GET("/status/:domain", Status)
	router.GET("/run/:runDomain", Run)
	router.GET("/stop/:stopDomain", Stop)
	env := os.Getenv("APP_ENV")
	if env == "prod" {
		log.Println("Running api server in prod mode")
	} else {
		log.Println("Running api server in dev mode")
	}
	http.ListenAndServe(":3030", router)
}
