// +heroku goVersion go1.14
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

const name = "urlín"
const version = "0.2.0"

var stdout, stderr io.Writer = os.Stdout, os.Stderr
var exit func(int) = os.Exit
var defaultPort string = "9000"

func parse(args ...string) (o struct {
	ba, ru, rc, wr string
	h, v, x        bool
}) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)

	fs.StringVar(&o.ba, "b", getEnv("BIND_ADDR", ":"+defaultPort), "address to which to bind")
	fs.StringVar(&o.ru, "r", getEnv("REDIS_URL", "redis://localhost:6379"), "redis host")
	fs.StringVar(&o.rc, "C", getEnv("REDIS_COUNTER", "_"), "redis counter")
	fs.StringVar(&o.wr, "w", getEnv("WEB_ROOT", "public"), "serve static files from this dir")
	fs.BoolVar(&o.x, "X", getEnv("ALLOW_CORS", "") == "true", "allow cross-origin requests")
	fs.BoolVar(&o.v, "v", false, "print the version number and exit")
	fs.BoolVar(&o.h, "h", false, "print help and exit")

	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}

	if o.v {
		fmt.Fprintf(stdout, "%s v%s\n", name, version)
	}

	if o.h {
		fs.SetOutput(stdout)
		fs.PrintDefaults()
	}

	return
}

func main() {
	opts := parse(os.Args[1:]...)
	if opts.h || opts.v {
		return
	}

	r := &redisStore{
		url:     opts.ru,
		counter: opts.rc,
	}

	if err := r.ping(); err != nil {
		log.Fatalf("failed to open redis at %s: %s", opts.ru, err)
	}

	h := &handler{store: r}

	if dirExists(opts.wr) {
		h.www = http.FileServer(http.Dir(opts.wr))
	}

	x := cors.Default()
	if opts.x {
		x = cors.AllowAll()
	}

	http.Handle("/", x.Handler(h))
	fmt.Printf("urleen listening on %s\n", opts.ba)
	if err := http.ListenAndServe(opts.ba, nil); err != nil {
		log.Fatal(err)
	}
}

func getEnv(name, fallback string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return fallback
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true
	}
	if !os.IsNotExist(err) {
		log.Printf("%s is not a readable directrory", path)
	}
	return false
}

func init() {
	if port, found := os.LookupEnv("PORT"); found {
		defaultPort = port
	}
}
