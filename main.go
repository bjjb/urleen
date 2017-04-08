package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/go-redis/redis"
)

const (
	Base62          = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultBindAddr = ":8089"
	DefaultWebRoot  = "www"
	DefaultRedisURL = "redis://localhost:6379"
)

var (
	bindAddr, webRoot, redisURL string
	base62                      = map[rune]int{}
	pattern                     *regexp.Regexp
	static                      http.Handler
	redisClient                 *redis.Client
)

func init() {
	for i, r := range Base62 {
		base62[r] = i
	}

	pattern = regexp.MustCompile("^/[a-zA-Z0-9]+$")

	var found bool
	if bindAddr, found = os.LookupEnv("BIND_ADDR"); !found {
		bindAddr = DefaultBindAddr
	}
	if webRoot, found = os.LookupEnv("WEB_ROOT"); !found {
		webRoot = DefaultWebRoot
	}
	if redisURL, found = os.LookupEnv("REDIS_URL"); !found {
		redisURL = DefaultRedisURL
	}

	flag.StringVar(&bindAddr, "b", bindAddr, "address/port to which to bind")
	flag.StringVar(&redisURL, "r", redisURL, "redis host")
	flag.StringVar(&webRoot, "w", webRoot, "directory holding static files")

}

func main() {
	flag.Parse()

	static = http.FileServer(http.Dir(webRoot))

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	redisClient = redis.NewClient(options)

	http.HandleFunc("/", handler)

	fmt.Printf("urleen listening on %s\n", bindAddr)
	if err := http.ListenAndServe(bindAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func redisOptions() *redis.Options {
	options := &redis.Options{}
	uri, err := url.Parse(redisURL)
	if err != nil {
		log.Fatalf("couldn't parse Redis URL %q; %q", redisURL, err)
	}
	options.Addr = uri.Host
	if pw, set := uri.User.Password(); set {
		options.Password = pw
	}
	if regexp.MustCompile(`^/\d+$`).MatchString(uri.Path) {
		if db, err := strconv.Atoi(uri.Path[1:]); err == nil {
			options.DB = db
		}
	}
	return options
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	switch r.Method {
	case http.MethodGet:
		switch {
		case pattern.MatchString(r.URL.Path):
			if location, found := getURL(r.URL.Path[1:]); found {
				http.Redirect(w, r, location, http.StatusMovedPermanently)
				return
			}
			http.NotFound(w, r)
		default:
			static.ServeHTTP(w, r)
		}
	case http.MethodPost:
		var location string
		if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id := putURL(location)
		if err := json.NewEncoder(w).Encode(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getURL(id string) (string, bool) {
	location, err := redisClient.Get(id).Result()
	if err == redis.Nil {
		return "", false
	}
	return location, true
}

func putURL(location string) string {
	id := encode62(uint64(redisClient.Incr("_").Val()))
	redisClient.Set(id, location, 0)
	return id
}

func decode62(s string) uint64 {
	var x uint64
	max := len(s) - 1
	for i, c := range s {
		x = x + uint64(base62[c]*int(math.Pow(float64(62), float64(max-i))))
	}
	return x
}

func encode62(i uint64) string {
	if i == 0 {
		return "0"
	}
	r := uint64(62)
	b := []byte{}
	for i > 0 {
		b = append([]byte{Base62[i%r]}, b...)
		i /= r
	}
	return string(b)
}
