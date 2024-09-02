package main

import (
	"github.com/jmoiron/sqlx"
)

type site struct {
	DB            *sqlx.DB
	BaseURL       string
	ImagesPerPage int
}

func newSite(db *sqlx.DB, base string) *site {
	return &site{
		DB:      db,
		BaseURL: base,
	}
}
