package main // import "github.com/thraxil/sofi

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Image struct {
	Id         int       `db:"id"`
	AHash      string    `db:"ahash"`
	Ext        string    `db:"ext"`
	Url        string    `db:"url"`
	InsertedAt time.Time `db:"inserted_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Tag struct {
	Id         int       `db:"id"`
	Slug       string    `db:"slug"`
	Tag        string    `db:"tag"`
	InsertedAt time.Time `db:"inserted_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func count_images(db *sqlx.DB) int {
	var count int
	db.Get(&count, "SELECT COUNT(*) FROM images")
	return count
}

func newest_images(db *sqlx.DB, page int) []Image {
	// page size is 20
	images := []Image{}
	db.Select(&images, "SELECT * FROM images ORDER BY id DESC LIMIT 20 OFFSET $1", (page-1)*20)
	return images
}

func random_image(db *sqlx.DB) Image {
	image := Image{}
	db.Get(&image, "SELECT * FROM images ORDER BY RANDOM() LIMIT 1")
	return image
}

var templateDir = "templates"

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

func makeHandler(fn func(http.ResponseWriter, *http.Request, *site), s *site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, s)
	}
}

type indexResponse struct {
	Images []Image
	Page   int
}

func indexHandler(w http.ResponseWriter, r *http.Request, s *site) {
	ir := indexResponse{}
	images := newest_images(s.DB, 1)
	ir.Images = images
	ir.Page = 1
	tmpl := getTemplate("index.html")
	tmpl.Execute(w, ir)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	// just ignore this crap
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(templateDir, "base.html"),
		filepath.Join(templateDir, filename),
	))
}

func healthzHandler(w http.ResponseWriter, _ *http.Request, _ *site) {
	w.WriteHeader(http.StatusOK)
}

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
	mux.HandleFunc("/", makeHandler(indexHandler, s))
	mux.HandleFunc("/favicon.ico", faviconHandler)
	http.ListenAndServe(":"+port, mux)
}
