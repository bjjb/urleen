package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	urlList := &urlList{[]byte("http://other.com")}

	post := func(u string) *http.Response {
		w := httptest.NewRecorder()
		body := strings.NewReader(fmt.Sprintf("%q", u))
		r := httptest.NewRequest("POST", "http://example.com/", body)
		r.Header.Set("Content-Type", "application/json")
		urlList.ServeHTTP(w, r)
		return w.Result()
	}

	get := func(path string) *http.Response {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("http://example.com%s", path), nil)
		urlList.ServeHTTP(w, r)
		return w.Result()
	}

	resp := get("/0")
	if resp.StatusCode != http.StatusPermanentRedirect {
		t.Errorf("expected a redirect, got %d (%s)", resp.StatusCode, resp.Status)
	}
	if resp.Header.Get("Location") != "http://other.com" {
		t.Errorf("expected a redirect to http://other.com, got %s", resp.Header.Get("Location"))
	}

	resp = get("/1")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected a 404, got %d (%s)", resp.StatusCode, resp.Status)
	}

	resp = post("http://foo.bar")
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected a 203, got %d (%s)", resp.StatusCode, resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "1" {
		t.Errorf("expected 1, got %s", string(body))
	}

	resp = get("/1")
	if resp.StatusCode != http.StatusPermanentRedirect {
		t.Errorf("expected a redirect, got %d (%s)", resp.StatusCode, resp.Status)
	}
	if resp.Header.Get("Location") != "http://foo.bar" {
		t.Errorf("expected a redirect to http://foo.bar, got %s", resp.Header.Get("Location"))
	}

	w := httptest.NewRecorder()
	urlList.ServeHTTP(w, httptest.NewRequest("POST", "http://example.com/", nil))
	resp = w.Result()
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("expected a 415, got %d (%s)", resp.StatusCode, resp.Status)
	}

	resp = post("nonsense")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400, got %d (%s)", resp.StatusCode, resp.Status)
	}

	w = httptest.NewRecorder()
	urlList.ServeHTTP(w, httptest.NewRequest("DELETE", "http://example.com/", nil))
	resp = w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected a 405, got %d (%s)", resp.StatusCode, resp.Status)
	}

	(*urlList)[0] = nil
	resp = get("/0")
	if resp.StatusCode != http.StatusGone {
		t.Errorf("expected a 410, got %d (%s)", resp.StatusCode, resp.Status)
	}

	resp = get("/")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200, got %d (%s)", resp.StatusCode, resp.Status)
	}
	contentType := resp.Header.Get("Content-Type")
	if !regexp.MustCompile("^text/html").MatchString(contentType) {
		t.Errorf("expected text/html, got %s", resp.Header.Get("Content-Type"))
	}

	w = httptest.NewRecorder()
	urlList.ServeHTTP(w, httptest.NewRequest("GET", "http://example.com/js/app.js", nil))
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200, got %d (%s)", resp.StatusCode, resp.Status)
	}
	contentType = resp.Header.Get("Content-Type")
	if !regexp.MustCompile("^application/javascript").MatchString(contentType) {
		t.Errorf("expected text/html, got %s", resp.Header.Get("Content-Type"))
	}

	w = httptest.NewRecorder()
	urlList.ServeHTTP(w, httptest.NewRequest("GET", "http://example.com/css/style.css", nil))
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200, got %d (%s)", resp.StatusCode, resp.Status)
	}
	contentType = resp.Header.Get("Content-Type")
	if !regexp.MustCompile("^text/css").MatchString(contentType) {
		t.Errorf("expected text/html, got %s", resp.Header.Get("Content-Type"))
	}
}
