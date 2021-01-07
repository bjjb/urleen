package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type handler struct {
	www   http.Handler
	store store
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveHTTP(w, r)
}

func (h *handler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet || r.Method == http.MethodHead:
		if h.store != nil && pattern.MatchString(r.URL.Path) {
			if location := h.store.get(r.URL.Path[1:]); location != "" {
				http.Redirect(w, r, location, http.StatusMovedPermanently)
				return
			}
			http.NotFound(w, r)
			return
		}
		if h.www != nil {
			h.www.ServeHTTP(w, r)
			return
		}
	case r.Method == http.MethodPost && h.store != nil:
		var location string
		if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id := h.store.put("", location)
		if err := json.NewEncoder(w).Encode(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		text := http.StatusText(http.StatusMethodNotAllowed)
		http.Error(w, text, http.StatusMethodNotAllowed)
	}
}

var pattern = regexp.MustCompile("^/[a-zA-Z0-9]+$")
