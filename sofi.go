package main // import "github.com/thraxil/sofi

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	baseUrl, exists := os.LookupEnv("SOFI_BASE_URL")
	if !exists {
		baseUrl = "http://localhost:8080"
	}
	port, exists := os.LookupEnv("SOFI_PORT")
	if !exists {
		port = "8080"
	}
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=melo_dev sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	s := newSite(db, baseUrl)

	mux := http.NewServeMux()
	mux.HandleFunc("/", makeHandler(indexHandler, s))
	mux.HandleFunc("/image/{id}", makeHandler(imageHandler, s))
	mux.HandleFunc("/random", makeHandler(randomHandler, s))
	mux.HandleFunc("/tag", makeHandler(tagIndexHandler, s))
	mux.HandleFunc("/tag/{tag}", makeHandler(tagHandler, s))
	mux.HandleFunc("/favicon.ico", faviconHandler)
	http.ListenAndServe(":"+port, mux)
}
