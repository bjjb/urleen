package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	radix    = len(alphabet)
	file     = "urls.gz"
)

var (
	addr   = ":8080"
	urls   = []string{}
	values = map[rune]int{}
	re     = regexp.MustCompile(fmt.Sprintf(`/[%s]+`, alphabet))
)

func init() {
	re = regexp.MustCompile(fmt.Sprintf(`/[%s]+`, alphabet))
	for i, r := range alphabet {
		values[r] = i
	}
	flag.StringVar(&addr, "a", ":8080", "specify the server address")
	http.HandleFunc("/", handle)
}

func main() {
	load()
	flag.Parse()
	go gracefulExit()
	if err := http.ListenAndServe(fmt.Sprintf("%s", addr), nil); err != nil {
		log.Fatal(err)
	}
}

func gracefulExit() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch // blocks
	save()
	signal.Reset(os.Interrupt)
	os.Exit(0)
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPost:
		post(w, r)
	default:
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	if re.MatchString(r.URL.Path) {
		url := urls[decode(r.URL.Path[1:])]
		if url == "" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, url, http.StatusMovedPermanently)
		return
	}
	http.NotFound(w, r)
}

func decode(s string) uint64 {
	var n uint64
	max := len(s) - 1
	for i, c := range s {
		n = n + uint64(values[c]*int(math.Pow(float64(radix), float64(max-i))))
	}
	return n
}

func encode(i uint64) string {
	if i == 0 {
		return "0"
	}
	r := uint64(radix)
	b := []byte{}
	for i > 0 {
		b = append([]byte{alphabet[i%r]}, b...)
		i /= r
	}
	return string(b)
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing request body; %v")
		}
	}()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body; %v", err)
		http.Error(w, "Uh oh", http.StatusInternalServerError)
		return
	}
	u := string(bytes.TrimSpace(data))
	if _, err := url.Parse(u); err != nil {
		http.Error(w, "Malformed request", http.StatusBadRequest)
		return
	}
	urls = append(urls, u)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", encode(uint64(len(urls)-1)))
}

func load() {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		urls = []string{"https://bjjb.github.io", "https://ddg.gg"}
		return
	}
	f, err := os.Open(file)
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("error closing file; %v", err)
		}
	}()
	if err != nil {
		log.Fatal("couldn't open file; %v", err)
	}
	zr, err := gzip.NewReader(f)
	defer func() {
		if err := zr.Close(); err != nil {
			log.Printf("error closing gzip reader; %v", err)
		}
	}()
	if err != nil {
		log.Fatalf("couldn't create zip reader; %v", err)
	}
	data, err := ioutil.ReadAll(zr)
	if err != nil {
		log.Fatalf("error reading from zip stream; %v", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			urls = append(urls, line)
		}
	}
}

func save() {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("error closing file; %v", err)
		}
	}()
	zw := gzip.NewWriter(f)
	defer func() {
		if err := zw.Close(); err != nil {
			log.Printf("error closing gzip writer; %v", err)
		}
	}()
	w := bufio.NewWriter(zw)
	for _, l := range urls {
		fmt.Fprintf(w, "%s\n", l)
	}
	if err := w.Flush(); err != nil {
		log.Printf("error flushing buffered writer; %v", err)
	}
}
