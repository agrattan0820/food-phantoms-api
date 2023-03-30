package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

	rows, err := s.DB.Query("SELECT k.*, c.name AS \"parent_name\", c.website_link AS \"parent_link\" FROM kitchens k LEFT JOIN companies c ON k.parent_id = c.id")
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var kitchen Kitchen
		if err := rows.Scan(&kitchen.ID, &kitchen.CreatedAt, &kitchen.UpdatedAt, &kitchen.Name, &kitchen.Logo, &kitchen.Description, &kitchen.WebsiteLink, &kitchen.ParentID, &kitchen.Type, &kitchen.Slug, &kitchen.DoorDashLink, &kitchen.ParentName, &kitchen.ParentLink); err != nil {
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

func (s *Server) KitchenBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var kitchens []Kitchen
	var kitchen Kitchen

	var locations []Location

	var companiesKitchenRunsIn []Company

	// Query for kitchen and its parent
	row := s.DB.QueryRow("SELECT k.*, c.name AS \"parent_name\", c.website_link AS \"parent_link\" FROM kitchens k LEFT JOIN companies c ON k.parent_id = c.id WHERE k.slug = $1", slug)

	switch err := row.Scan(&kitchen.ID, &kitchen.CreatedAt, &kitchen.UpdatedAt, &kitchen.Name, &kitchen.Logo, &kitchen.Description, &kitchen.WebsiteLink, &kitchen.ParentID, &kitchen.Type, &kitchen.Slug, &kitchen.DoorDashLink, &kitchen.ParentName, &kitchen.ParentLink); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		w.WriteHeader(http.StatusNotFound)
		return
	case nil:
		fmt.Println(kitchen)
		kitchens = append(kitchens, kitchen)
	default:
		panic(err)
	}

	// Query for locations
	rows, err := s.DB.Query("SELECT * FROM locations WHERE kitchen_id = $1", kitchen.ID)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var location Location
		if err := rows.Scan(&location.ID, &location.CreatedAt, &location.UpdatedAt, &location.KitchenID, &location.Address1, &location.City, &location.State, &location.Country, &location.ZipCode, &location.GoogleRating, &location.Address2); err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		locations = append(locations, location)
	}

	// Query for runs-in companies
	rows, err = s.DB.Query("SELECT c.* FROM kitchen_runs_in_company kc JOIN companies c ON kc.company_id = c.id WHERE kitchen_id = $1", kitchen.ID)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var company Company
		if err := rows.Scan(&company.ID, &company.CreatedAt, &company.Name, &company.Description, &company.Logo, &company.WebsiteLink, &company.UpdatedAt); err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		companiesKitchenRunsIn = append(companiesKitchenRunsIn, company)
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setupPayload := KitchenByIdPayload{
		Kitchen:                kitchens,
		Locations:              locations,
		CompaniesKitchenRunsIn: companiesKitchenRunsIn,
	}

	payload, err := json.Marshal(setupPayload)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload = TrimKitchen(payload)

	w.Write(payload)
}

type AddKitchenRequest struct {
	Name         *string `json:"name"`
	DoorDashLink *string `json:"doordash_link"`
	WebsiteLink  *string `json:"website_link"`
	Parent       *string `json:"parent"`
}

func (s *Server) AddKitchen(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	data := &AddKitchenRequest{}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Failed to unmarshal payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sqlStatement := `
INSERT INTO kitchen_requests (kitchen_name, doordash_link, website_link, parent)
VALUES ($1, $2, $3, $4)
RETURNING id`
	id := 0
	err = s.DB.QueryRow(sqlStatement, data.Name, data.DoorDashLink, data.WebsiteLink).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New Kitchen Request ID is:", id)

	w.Write([]byte("Success"))
}
