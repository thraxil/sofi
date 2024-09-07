package main

import (
	"net/http"
)

func addRoutes(mux *http.ServeMux, s *site) {
	mux.HandleFunc("/", makeHandler(indexHandler, s))
	mux.HandleFunc("/image/{id}", makeHandler(imageHandler, s))
	mux.HandleFunc("/random", makeHandler(randomHandler, s))
	mux.HandleFunc("/tag", makeHandler(tagIndexHandler, s))
	mux.HandleFunc("/tag/{tag}", makeHandler(tagHandler, s))
	mux.HandleFunc("/feeds/newest", makeHandler(feedHandler, s))
	mux.HandleFunc("/smoketest/", makeHandler(healthzHandler, s))
	mux.HandleFunc("/favicon.ico", faviconHandler)
}
