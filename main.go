package main

import (
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
		r.ParseForm()
		x := r.Form.Get("url")
		uri, err := url.Parse(x)
		if err != nil || uri.Host == "" && uri.Scheme == "" {
			http.Error(w, "URL is not valid", http.StatusBadRequest)
			return
		}
		if len(*u) > (1<<32 - 2) {
			http.Error(w, "URL list is full", http.StatusInsufficientStorage)
			return
		}
		id := base62.Encode(uint64(len(*u)))
		*u = append(*u, []byte(x))
		w.WriteHeader(http.StatusCreated)
		location := url.URL{Scheme: r.URL.Scheme, Opaque: r.URL.Opaque, User: r.URL.User, Host: r.URL.Host, Path: id}
		fmt.Fprintf(w, "%s", location.String())
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
			return
		}
		staticHandler.ServeHTTP(w, r)
	default:
		http.Error(w, "GET short or POST long", http.StatusMethodNotAllowed)
	}
}
