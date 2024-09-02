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
	Images      []Image
	Page        int
	TotalPages  int
	NextPage    int
	HasNextPage bool
	PrevPage    int
	HasPrevPage bool
}

func indexHandler(w http.ResponseWriter, r *http.Request, s *site) {
	imagesPerPage := 20
	pagen, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		// can't parse as int? just default to one
		pagen = 1
	}
	ir := indexResponse{}
	cnt_images := count_images(s.DB)
	has_next_page := true
	has_prev_page := true
	if pagen < 1 {
		pagen = 1
		has_prev_page = false
	}
	if pagen > cnt_images/imagesPerPage {
		pagen = cnt_images / imagesPerPage
		has_next_page = false
	}
	images := newest_images(s.DB, pagen)
	ir.Images = images
	ir.Page = pagen
	ir.TotalPages = cnt_images / imagesPerPage
	ir.NextPage = pagen + 1
	ir.HasNextPage = has_next_page
	ir.PrevPage = pagen - 1
	ir.HasPrevPage = has_prev_page
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

func tagHandler(w http.ResponseWriter, r *http.Request, s *site) {
}

type tagIndexResponse struct {
	Tags []Tag
}

func tagIndexHandler(w http.ResponseWriter, r *http.Request, s *site) {
	tir := tagIndexResponse{}
	tags := getAllTags(s.DB)
	tir.Tags = tags
	tmpl := getTemplate("tags.html")
	tmpl.Execute(w, tir)
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
