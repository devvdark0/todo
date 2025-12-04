package main

import (
	"log"
	"net/http"

	"github.com/devvdark0/todo/pkg/db"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("successfully connected to mariadb!")
	_ = database
	log.Println("starting server on port :80...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
