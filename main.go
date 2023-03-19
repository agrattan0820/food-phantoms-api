package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var stats = map[string]int{}
var calls []string

type server struct {
	db *sql.DB
}

type Kitchen struct {
	ID          int8
	CreatedAt   string
	UpdatedAt   string
	Name        string
	Logo        sql.NullString
	Description sql.NullString
	WebsiteLink sql.NullString
	ParentID    sql.NullInt16
	Type        string
}

func (s *server) hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	calls = append(calls, name)
	stats[name]++

	fmt.Printf("calls: %#v\n", calls)
	fmt.Printf("stats: %#v\n\n", stats)

	fmt.Fprint(w, "Hello, ", name)
}

func (s *server) kitchens(w http.ResponseWriter, r *http.Request) {

	var kitchens []Kitchen

	rows, err := s.db.Query("SELECT * FROM kitchens")
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var kitchen Kitchen
		if err := rows.Scan(&kitchen.ID, &kitchen.CreatedAt, &kitchen.UpdatedAt, &kitchen.Name, &kitchen.Logo, &kitchen.Description, &kitchen.WebsiteLink, &kitchen.ParentID, &kitchen.Type); err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		kitchens = append(kitchens, kitchen)
	}

	if kitchens == nil {
		w.WriteHeader(404)
		w.Write([]byte("kitchens not found"))
		return
	}

	payload, err := json.Marshal(kitchens)

	if err != nil {
		log.Println("Failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(payload)

}

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
	s := server{db: db}

	// Test db by sending a ping
	pingErr := s.db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Starting server on :8080...")
	fmt.Println("Connected!")

	// Set up chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Start server
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/hello", s.hello)
	r.Get("/kitchens", s.kitchens)
	log.Fatal(http.ListenAndServe(":8080", r))
}
