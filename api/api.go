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
	"labix.org/v2/mgo/bson"
)

var texts *mgo.Collection

type Text struct {
	Content string
	URL     string
	Source  string
	IsHS    bool
	IsNotHS bool
	Idk     bool
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	session, err := database.Init()
	defer session.Close()
	texts = session.DB("app").C("texts")

	result := []Text{}
	err = texts.Find(nil).All(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	responseJSON(w, result)
}

func Random(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	session, err := database.Init()
	defer session.Close()
	texts = session.DB("app").C("texts")

	pipe := texts.Pipe([]bson.M{{"$sample": bson.M{"size": 1}}})
	result := []bson.M{}
	err = pipe.All(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	responseJSON(w, result)
}

func Submit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	session, _ := database.Init()
	defer session.Close()
	texts = session.DB("app").C("texts")

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

	if err := texts.Insert(text); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, text)
}

func Remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	session, err := database.Init()
	defer session.Close()
	texts = session.DB("app").C("texts")

	id := ps.ByName("id")
	err = texts.Remove(bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	responseJSON(w, "")
}

// used for COR preflight checks
func corsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
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

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func main() {
	//	session, err := database.Init()
	//	defer session.Close()
	//	texts = session.DB("app").C("texts")
	//	if err != nil {
	//		log.Println("connection to mongodb failed, aborting...")
	//		log.Fatal(err)
	//	}

	log.Println("connected to mongodb")

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/random", Random)
	router.POST("/submit", Submit)
	router.DELETE("/remove/:id", Remove)

	env := os.Getenv("APP_ENV")
	if env == "prod" {
		log.Println("Running api server in prod mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	http.ListenAndServe(":5000", router)
}
