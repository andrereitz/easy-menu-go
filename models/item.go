package models

import "database/sql"

type ItemData struct {
	Id          int            `json:"id"`
	Category    sql.NullInt64  `json:"category"`
	User        int            `json:"user"`
	MediaId     sql.NullInt64  `json:"media_id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Price       float64        `json:"price"`
}
