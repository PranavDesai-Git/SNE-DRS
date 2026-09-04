package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/mattn/go-sqlite3" // Important! Registers the sqlite3 driver
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to connect to db %e", err)
	}
	defer db.Close()

	if len(os.Args) > 1 && os.Args[1] == "import" {
		RunImport(db)
		return
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	r.Get("/habitations", getHabitations)
	r.Get("/risk-zones", getRiskZones)
	r.Get("/sites", getSites)

	// Serve the compiled Svelte frontend
	fs := http.FileServer(http.Dir("../frontend/dist"))
	r.Handle("/*", http.StripPrefix("/", fs))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
