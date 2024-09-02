package main // import "github.com/thraxil/sofi

import (
	"fmt"
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

	fmt.Println(count_images(db))
	fmt.Println(random_image(db))
	fmt.Println(newest_images(db, 1))

	mux := http.NewServeMux()
	mux.HandleFunc("/image/{id}", makeHandler(imageHandler, s))
	mux.HandleFunc("/", makeHandler(indexHandler, s))
	mux.HandleFunc("/random", makeHandler(randomHandler, s))
	mux.HandleFunc("/favicon.ico", faviconHandler)
	http.ListenAndServe(":"+port, mux)
}
