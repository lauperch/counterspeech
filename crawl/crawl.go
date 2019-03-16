package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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
		go scrape(src)
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
		// TODO get channel to src, send msg to src to stop
		src := sources[domain]
		src.isRunning = false
		sources[domain] = src
		responseJSON(w, "stopped scraping "+domain)
	}
}

func scrape(src Source) {
	startUrlHtml, err := getHtml(src.startUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	commentLinks := getCommentLinks(startUrlHtml)
	c := make(chan Text)
	for _, link := range commentLinks {
		go htmlToText(link, c)
	}

	for text := range c {
		go save(text)
	}
}

func getCommentLinks(html string) []string {
	re := regexp.MustCompile(`h3.*data-vr-contentbox.*\/.*"`)
	linkStrings := re.FindAllString(html, -1)
	var links []string
	for _, linkString := range linkStrings {
		parts := strings.Split(linkString, " ")
		if len(parts) == 6 {
			link := strings.Replace(parts[3], `"><a`, "", 1)
			links = append(links, "https://www.20min.ch"+link)
		}
	}
	return links
}

func htmlToText(url string, c chan Text) {
	re := regexp.MustCompile(`<p class="content">.*<`)
	html, err := getHtml(url)
	if err != nil {
		log.Print(err.Error())
		return
	}
	commentStrings := re.FindAllString(html, -1)
	for _, commentString := range commentStrings {
		parts := strings.Split(commentString, ">")
		if len(parts) == 2 {
			comment := strings.Replace(parts[1], " <", "", -1)
			comment = fmt.Sprintf("%q", comment)
			comment = strings.Replace(comment, "\\xe4", "ä", -1)
			comment = strings.Replace(comment, "\\xf6", "ö", -1)
			comment = strings.Replace(comment, "\\xfc", "ü", -1)
			comment = strings.Replace(comment, `\"`, `"`, -1)
			comment = strings.Replace(comment, `"`, "", -1)
			text := Text{Content: comment, URL: url}
			c <- text
		}
	}
}

func save(text Text) {
	url := ""
	if os.Getenv("APP_ENV") == "prod" {
		url = "http://35.198.123.101:5000/submit"
	} else {
		url = "http://192.168.0.67:5000/submit"
	}
	textJson, _ := json.Marshal(text)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(textJson))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getHtml(url string) (string, error) {
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
