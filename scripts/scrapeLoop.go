package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting scrape loop...")
	run("https://www.20min.ch/schweiz/")
	run("https://www.20min.ch/wissen/")
	run("https://www.20min.ch/digital/")
	run("https://www.20min.ch/entertainment/")
	run("https://www.20min.ch/people/")
}

func run(startUrl string) {
	r, err := http.Get("http://35.198.123.101:3030/run/www.20min.ch?startUrl=" + startUrl)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(r)
	r, err = http.Get("http://35.198.123.101:3030/stop/www.20min.ch")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(r)
}
