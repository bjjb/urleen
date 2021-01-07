package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/bjjb/urleen/base62"
	"github.com/go-redis/redis/v8"
)

type store interface {
	put(string) string
	get(string) string
}

type redisStore struct {
	client  *redis.Client
	counter string
}

func (r *redisStore) put(s string) string {
	ctx := context.Background()
	id := base62.Encode(uint64(r.client.Incr(ctx, r.counter).Val()))
	r.client.Set(ctx, id, s, 0)
	return id
}

func (r *redisStore) get(id string) string {
	ctx := context.Background()
	s, _ := r.client.Get(ctx, id).Result()
	return s
}

type mapStore struct {
	counter uint64
	data    map[string]string
	mutex   *sync.Mutex
}

func (m *mapStore) put(s string) string {
	m.mutex.Lock()
	id := base62.Encode(m.counter)
	m.counter++
	m.mutex.Unlock()
	m.data[id] = s
	return id
}

func (m *mapStore) get(id string) string {
	return m.data[id]
}

type handler struct {
	www   http.Handler
	store store
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		id := h.store.put(location)
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

func main() {
	wr := getEnv("WEB_ROOT", "www")
	ba := getEnv("BIND_ADDR", ":9000")
	ru := getEnv("REDIS_URL", "redis://localhost:6379")
	rc := getEnv("REDIS_COUNTER", "_")

	flag.StringVar(&ba, "b", ba, "address to which to bind")
	flag.StringVar(&ru, "r", ru, "redis host")
	flag.StringVar(&rc, "C", rc, "redis counter")
	flag.StringVar(&wr, "w", wr, "directory holding the app")
	flag.Parse()

	h := &handler{}

	if fileExists(wr) {
		h.www = http.FileServer(http.Dir(wr))
	} else {
		log.Printf("no such file or directory: %s (api only)", wr)
	}

	if options, err := redis.ParseURL(ru); err == nil {
		h.store = &redisStore{
			client:  redis.NewClient(options),
			counter: rc,
		}
	} else {
		log.Printf("failed to open redis at %s: %s", ru, err)
	}

	if h.store != nil || h.www != nil {
		http.Handle("/", h)
		fmt.Printf("urleen listening on %s\n", ba)
		if err := http.ListenAndServe(ba, nil); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Print("nothing to serve, quitting.")
	}
}

func getEnv(name, fallback string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return fallback
}

func fileExists(path string) bool {
	info, _ := os.Stat(path)
	return info != nil && info.IsDir()
}
