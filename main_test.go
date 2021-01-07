package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_parse(t *testing.T) {
	out := new(bytes.Buffer)
	defer func(w io.Writer) { stdout = w }(stdout)
	stdout = out
	opts := parse("-v")
	want := fmt.Sprintf("%s v%s", name, version)
	got := strings.TrimSpace(out.String())
	if !regexp.MustCompile(want).MatchString(got) {
		t.Errorf("expected %q, got %q", want, got)
	}
	if !opts.v {
		t.Errorf("expected opts.v to be set")
	}

	opts = parse("-h")
	want = fmt.Sprintf("-h")
	got = strings.TrimSpace(out.String())
	if !regexp.MustCompile(want).MatchString(got) {
		t.Errorf("expected %q, got %q", want, got)
	}
	if !opts.h {
		t.Errorf("expected opts.h to be set")
	}
}

func Test_main(t *testing.T) {
	defer func(a []string) { os.Args = a }(os.Args)
	os.Args = []string{""}
	go main()
	time.Sleep(5 * time.Millisecond)
}
