package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting scrape loop...")
	go run("https://www.20min.ch/schweiz/")
}

func run(startUrl string) {
	r, err := http.Get("https://ddd/run/www.20min.ch?startUrl=" + startUrl)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
