package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

var templateDir = "templates"

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

type randomResponse struct {
	Image Image
	Tags  []Tag
}

func randomHandler(w http.ResponseWriter, r *http.Request, s *site) {
	rr := randomResponse{}
	image := random_image(s.DB)
	rr.Image = image
	tmpl := getTemplate("random.html")
	tmpl.Execute(w, rr)
}

type imageResponse struct {
	Image Image
	Tags  []Tag
}

func imageHandler(w http.ResponseWriter, r *http.Request, s *site) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}
	ir := randomResponse{}
	image := get_image_by_id(s.DB, id)
	ir.Image = image
	tmpl := getTemplate("image.html")
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
