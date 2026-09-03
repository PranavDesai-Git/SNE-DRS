package main

import (
	"database/sql"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to connect to db %e", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	r.Get("/habitations", getHabitations)
	r.Get("/risk-zones", getRiskZones)
	r.Get("/sites", getSites)
}
