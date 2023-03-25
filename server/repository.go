package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Kitchen struct {
	ID          int8       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name"`
	Logo        NullString `json:"logo"`
	Description NullString `json:"description"`
	WebsiteLink NullString `json:"website_link"`
	ParentID    NullInt64  `json:"parent_id"`
	Type        string     `json:"type"`
	Slug        string     `json:"slug"`
	ParentName  NullString `json:"parent_name"`
	ParentLink  NullString `json:"parent_link"`
}

type Location struct {
	ID           int8       `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	KitchenID    int8       `json:"kitchen_id"`
	Address1     string     `json:"address_1"`
	Address2     NullString `json:"address_2"`
	City         string     `json:"city"`
	State        string     `json:"state"`
	Country      string     `json:"country"`
	ZipCode      string     `json:"zip_code"`
	GoogleRating string     `json:"google_rating"`
}
type Company struct {
	ID          int8       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name"`
	Description NullString `json:"description"`
	Logo        NullString `json:"logo"`
	WebsiteLink NullString `json:"website_link"`
}

type KitchenByIdPayload struct {
	Kitchen                []Kitchen  `json:"kitchen"`
	Locations              []Location `json:"locations"`
	CompaniesKitchenRunsIn []Company  `json:"runs_in"`
}

func TrimKitchen(b []byte) []byte {
	s := string(b)

	trimLeft := strings.Replace(s, "[", "", 1)
	t := strings.Replace(trimLeft, "]", "", 1)

	return []byte(t)
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
// func (ni *NullInt64) UnmarshalJSON(b []byte) error {
//  err := json.Unmarshal(b, &ni.Int64)
//  ni.Valid = (err == nil)
//  return err
// }

// NullBool is an alias for sql.NullBool data type
type NullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON for NullBool
// func (nb *NullBool) UnmarshalJSON(b []byte) error {
//  err := json.Unmarshal(b, &nb.Bool)
//  nb.Valid = (err == nil)
//  return err
// }

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON for NullFloat64
// func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
//  err := json.Unmarshal(b, &nf.Float64)
//  nf.Valid = (err == nil)
//  return err
// }

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
// func (ns *NullString) UnmarshalJSON(b []byte) error {
//  err := json.Unmarshal(b, &ns.String)
//  ns.Valid = (err == nil)
//  return err
// }

// NullTime is an alias for mysql.NullTime data type
type NullTime struct {
	sql.NullTime
}

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}
