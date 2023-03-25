package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"food-phantoms-api/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to db
	connectionString := os.Getenv("DATABASE_URL")
	fmt.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set up server
	s := server.Server{DB: db}

	// Test db by sending a ping
	pingErr := s.DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Starting server on :8080...")
	fmt.Println("Connected! http://localhost:8080")

	// Set up chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Start server
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Get("/kitchens", s.Kitchens)
	r.Get("/kitchen/{id}", s.KitchenById)
	log.Fatal(http.ListenAndServe(":8080", r))
}
