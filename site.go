package main

import (
	"github.com/jmoiron/sqlx"
)

type site struct {
	DB            *sqlx.DB
	ImagesPerPage int
}

func newSite(db *sqlx.DB) *site {
	return &site{
		DB: db,
	}
}
