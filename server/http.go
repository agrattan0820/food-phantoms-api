package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

type Server struct {
	DB *sql.DB
}

func (s *Server) Kitchens(w http.ResponseWriter, r *http.Request) {

	var kitchens []Kitchen

	rows, err := s.DB.Query("SELECT * FROM kitchens")
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

func (s *Server) KitchenById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var kitchens []Kitchen
	var kitchen Kitchen

	row := s.DB.QueryRow("SELECT * FROM kitchens WHERE id = $1", id)

	switch err := row.Scan(&kitchen.ID, &kitchen.CreatedAt, &kitchen.UpdatedAt, &kitchen.Name, &kitchen.Logo, &kitchen.Description, &kitchen.WebsiteLink, &kitchen.ParentID, &kitchen.Type); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(kitchen)
		kitchens = append(kitchens, kitchen)
	default:
		panic(err)
	}

	payload, err := json.Marshal(kitchens)

	if err != nil {
		log.Println("Failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(payload)
}
