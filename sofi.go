package main // import "github.com/thraxil/sofi

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	port, exists := os.LookupEnv("SOFI_PORT")
	if !exists {
		port = "8080"
	}
	templateDir, exists = os.LookupEnv("SOFI_TEMPLATE_DIR")
	if !exists {
		templateDir = "templates"
	}
	var DB_URL string
	if os.Getenv("DATABASE_URL") != "" {
		DB_URL = os.Getenv("DATABASE_URL")
	} else {
		// local dev settings
		DB_URL = "user=postgres password=postgres dbname=melo_dev sslmode=disable"
	}
	db, err := sqlx.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalln(err)
	}

	s := newSite(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/", makeHandler(indexHandler, s))
	mux.HandleFunc("/image/{id}", makeHandler(imageHandler, s))
	mux.HandleFunc("/random", makeHandler(randomHandler, s))
	mux.HandleFunc("/tag", makeHandler(tagIndexHandler, s))
	mux.HandleFunc("/tag/{tag}", makeHandler(tagHandler, s))
	mux.HandleFunc("/feeds/newest", makeHandler(feedHandler, s))
	mux.HandleFunc("/smoketest/", makeHandler(healthzHandler, s))
	mux.HandleFunc("/favicon.ico", faviconHandler)
	http.ListenAndServe(":"+port, mux)
}
