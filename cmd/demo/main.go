package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fortytw2/gdgberlin"
)

func main() {
	log.Println("howdy folks")

	db, err := gdgberlin.NewDB(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", gdgberlin.NewHandler(db)))
}
