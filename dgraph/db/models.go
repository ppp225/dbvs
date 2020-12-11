package db

import "time"

type Item struct {
	UID         string     `json:"uid"`
	Type        []string   `json:"dgraph.type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type ItemList struct {
	Items []Item `json:"items"`
}
