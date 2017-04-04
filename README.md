Urlín
=====

A little URL squisher, in Go.

Build it with `go build` and install it with `go install`. Get help with
`urleen -h`. If you don't like the front-end, use a different one, and either
let urleen serve it (with the `-w` flag) or put it behind a server; urleen
will respond to `GET /<id>` with a redirect, and `POST /` by remembering the
(JSON) contents of the request body, provided it's a valid URL.

There's a Dockerfile included, so you can run it anywhere.
