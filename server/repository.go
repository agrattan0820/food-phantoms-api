package server

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Kitchen struct {
	ID          int8           `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Name        string         `json:"name"`
	Logo        sql.NullString `json:"logo"`
	Description sql.NullString `json:"description"`
	WebsiteLink sql.NullString `json:"website_link"`
	ParentID    sql.NullInt16  `json:"parent_id"`
	Type        string         `json:"type"`
}
