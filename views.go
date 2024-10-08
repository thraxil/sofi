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
	cntImages := countImages(s.DB)
	hasNextPage := true
	hasPrevPage := true
	if pagen < 1 {
		pagen = 1
	}
	if pagen < 2 {
		hasPrevPage = false
	}
	if pagen > cntImages/imagesPerPage {
		pagen = cntImages / imagesPerPage
		hasNextPage = false
	}
	images := newestImages(s.DB, pagen)
	ir.Images = images
	ir.Page = pagen
	ir.TotalPages = cntImages / imagesPerPage
	ir.NextPage = pagen + 1
	ir.HasNextPage = hasNextPage
	ir.PrevPage = pagen - 1
	ir.HasPrevPage = hasPrevPage
	tmpl := getTemplate("index.html")
	tmpl.Execute(w, ir)
}

type randomResponse struct {
	Image Image
	Tags  []Tag
}

func randomHandler(w http.ResponseWriter, r *http.Request, s *site) {
	rr := randomResponse{}
	image := randomImage(s.DB)
	rr.Image = image
	rr.Tags = getImageTags(s.DB, image)
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
	image := getImageById(s.DB, id)
	ir.Image = image
	ir.Tags = getImageTags(s.DB, image)
	tmpl := getTemplate("image.html")
	tmpl.Execute(w, ir)
}

type tagResponse struct {
	Tag    Tag
	Images []Image
}

func tagHandler(w http.ResponseWriter, r *http.Request, s *site) {
	slug := r.PathValue("tag")
	tr := tagResponse{}
	tag := getTagBySlug(s.DB, slug)
	tr.Tag = tag
	images := getImagesByTag(s.DB, tag)
	tr.Images = images
	tmpl := getTemplate("tag.html")
	tmpl.Execute(w, tr)
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

type feedResponse struct {
	Images []Image
}

func feedHandler(w http.ResponseWriter, r *http.Request, s *site) {
	pagen := 1
	ir := feedResponse{}
	images := newestImages(s.DB, pagen)
	ir.Images = images
	var t = template.New("feed.html")
	tmpl := template.Must(t.ParseFiles(filepath.Join(templateDir, "feed.html")))
	w.Header().Set("Content-Type", "application/rss+xml")
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
