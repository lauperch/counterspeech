package main

import (
	"database"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

var texts *mgo.Collection

type Text struct {
	Content string
	URL     string
}

func submit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	text := &Text{}
	err = json.Unmarshal(data, text)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert new post
	if err := texts.Insert(text); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, text)

}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	result := []Text{}
	err := texts.Find(nil).All(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	responseJSON(w, result)
}

// used for COR preflight checks
func corsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
}

// util
func getFrontendUrl() string {
	if os.Getenv("APP_ENV") == "prod" {
		return "http://localhost:3000" // change this to prod domain
	} else {
		return "http://localhost:3000"
	}
}

func setCors(w http.ResponseWriter) {
	frontendUrl := getFrontendUrl()
	w.Header().Set("Access-Control-Allow-Origin", frontendUrl)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Temporary Canary test to make sure Travis-CI is working
func Canary(word string) string {
	return word
}

func main() {
	session, err := database.Init()
	defer session.Close()

	texts = session.DB("app").C("texts")

	router := httprouter.New()
	router.GET("/", index)
	router.POST("/submit", submit)

	if err != nil {
		log.Println("connection to mongodb failed, aborting...")
		log.Fatal(err)
	}

	log.Println("connected to mongodb")

	env := os.Getenv("APP_ENV")
	if env == "prod" {
		log.Println("Running api server in prod mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	http.ListenAndServe(":8080", router)
}
