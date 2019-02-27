package main

import (
	"database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type Text struct {
	Content string
	URL     string
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	mySession := database.Session.Copy()
	defer mySession.Close()

	c := mySession.DB("test").C("myCollection")

	err := c.Insert(
		&Text{"Lorem Ipsum", "http://a.com"},
		&Text{"Dolor", "http://b.com"})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	result := Text{}
	err = c.Find(bson.M{"url": "http://b.com"}).One(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, result.Content)
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

// Temporary Canary test to make sure Travis-CI is working
func Canary(word string) string {
	return word
}

func main() {
	// add database
	session, err := database.Init()
	defer session.Close()

	// add router and routes
	router := httprouter.New()
	router.GET("/", indexHandler)

	if err != nil {
		log.Println("connection to mongodb failed, aborting...")
		log.Fatal(err)
	}

	log.Println("connected to mongodb")

	// print env
	env := os.Getenv("APP_ENV")
	if env == "prod" {
		log.Println("Running api server in prod mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	http.ListenAndServe(":8080", router)
}
