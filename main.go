package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/bjjb/urleen/base62"
)

var (
	bindAddr, webRoot string
	staticHandler     http.Handler
	re                = regexp.MustCompile("^/[a-zA-Z0-9]+$")
)

type urlList [][]byte

func init() {
	flag.StringVar(&bindAddr, "b", ":8089", "address/port to which to bind")
	flag.StringVar(&webRoot, "w", "www", "directory from which to serve static files")
	staticHandler = http.FileServer(http.Dir(webRoot))
}

func main() {
	flag.Parse()
	http.Handle("/", &urlList{})
	fmt.Printf("urleen listening on %s\n", bindAddr)
	if err := http.ListenAndServe(bindAddr, nil); err != nil {
		log.Fatal(err)
	}
}

// ServeHTTP implements the Handler for a urlList.
func (u *urlList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		inputURI := ""
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "try application/json", http.StatusUnsupportedMediaType)
			return
		}
		if len(*u) > (1<<32 - 2) {
			http.Error(w, "URL list is full", http.StatusInsufficientStorage)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&inputURI); err != nil {
			log.Printf("error decoding body: %s", err)
			http.Error(w, "failed to decode JSON", http.StatusInternalServerError)
			return
		}
		uri, err := url.Parse(inputURI)
		if err != nil || uri.Host == "" && uri.Scheme == "" {
			http.Error(w, "URL is not valid", http.StatusBadRequest)
			return
		}
		id := base62.Encode(uint64(len(*u)))
		*u = append(*u, []byte(inputURI))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		log.Printf("%s → %s", id, uri.String())
		fmt.Fprintf(w, "%q", id)
	case r.Method == http.MethodGet:
		if re.MatchString(r.URL.Path) {
			i := int(base62.Decode(r.URL.Path[1:]))
			if i >= len(*u) {
				http.NotFound(w, r)
				return
			}
			url := (*u)[i]
			if len(url) == 0 {
				http.Error(w, "gone", http.StatusGone)
				return
			}
			http.Redirect(w, r, string(url), http.StatusPermanentRedirect)
			log.Printf("%s → %s", r.URL.String(), string(url))
			return
		}
		staticHandler.ServeHTTP(w, r)
	default:
		http.Error(w, "GET short or POST long", http.StatusMethodNotAllowed)
	}
}
