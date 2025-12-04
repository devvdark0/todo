package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("starting server on port :80...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
