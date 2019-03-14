package main

import (
	"bytes"
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

var pages = []string{"www.20min.ch"}

func Status(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	domain := ps.ByName("domain")
	if !contains(pages, domain) {
		responseJSON(w, "domain not yet implemented")
	} else if sources[domain].isRunning {
		responseJSON(w, "source "+domain+" is running")
	} else {
		responseJSON(w, "source "+domain+" is not running")
	}
}

func Run(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	domain := ps.ByName("runDomain")
	if !contains(pages, domain) {
		responseJSON(w, "domain not yet implemented")
	} else if sources[domain].isRunning {
		responseJSON(w, "source"+domain+" already running")
	} else {
		startUrl := r.URL.Query().Get("startUrl")
		src := Source{domain: domain, startUrl: startUrl}
		go Scrape(src)
		src.isRunning = true
		sources[domain] = src
		responseJSON(w, "started scraping "+domain+" on url "+startUrl)
	}
}

func Stop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	domain := ps.ByName("stopDomain")
	if !contains(pages, domain) {
		responseJSON(w, "domain not yet implemented")
	} else if !sources[domain].isRunning {
		responseJSON(w, "source"+domain+" already stopped")
	} else {
		src := sources[domain]
		src.isRunning = false
		sources[domain] = src
		responseJSON(w, "stopped scraping "+domain)
	}
}

func Scrape(src Source) {
	startUrlHtml, err := GetHtml(src.startUrl)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(startUrlHtml)
	}
}

func NextUrl(prevUrl string) string {
	// TODO implement
	return ""
}

func HtmlToText(html string) Text {
	// TODO implement
	return Text{}
}

func Save(text Text) {
	url := ""
	if os.Getenv("APP_ENV") == "prod" {
		url = "" // TODO insert correct prod url
	} else {
		url = "http://localhost:5000"
	}
	textJson, _ := json.Marshal(text)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(textJson))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetHtml(url string) (string, error) {
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
