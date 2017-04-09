URLín
=====

A little URL squisher, backed by Redis, a very small Go server, and with a
teensy clean HTML5 front-end. Run it on your secure short domain, and spam
the internets with pretty little URLs.

Here's a screenshot:

![A squished
URL](https://lh3.googleusercontent.com/WPoxS8zEYhP1gD48mchJm1G2GVa7uCzYpsyrMaH3w007wr6vvnH8TTZUTz3mmkgpCcW1NqcMGO5hvJHXYoCPjkkX1TUt_hZoHSyfxV5eZbZ9zjWmzNZTJsWvVvqebCEPFoy-mIXChDoLmMjsEBYlUmzXHsI4pD3o3yJZcok5XWD4JQHo4kaDkKAOf9nAtSsuGtwZeN7RhCvNj6hxud4HFvDZKswC1HkULDrPcfT8Ot-cwC5sJMXuym4ipcPKnLxWxY5yIykNrHSO6NaGIntin8bfTKf2lXQWCLrBc9AQsvJSjrzDAV3JJ2ZzawMAA-nu-1nR9xQt6-wxk8fzWhiyw3TvpyAoCVuz3atuiNPWVKEU9vHJgwJW1Fyx-CEVKn8ilRtZEPaDb9WxI3V05Cj3jvjJ_gT1X4zU0236XVXid47x3BpSdjsnLoCzkRIz7zXcThRovGDRQlvWbUzyYq1-_UaFImkTy9NID5oASZBjvnfD-rwMVxCCQZNTVVxGzuuvR5GO6kbyG__QqmZ_ugj55aU9pvsvNwVvLI2-Q_945WG7XlJqDQqB9qcyOkYCHze64jmUkl3EoFh71yOESins9A--qV4H8JIGxHqibOqFeuKgnwsmi1jzgFwJteqMXY4xcQoVMtJDZhyca9jPnbfe6oEsGXQ3npuYj2AOFdG9PTk=w385-h684-no "Isn't it pretty?")

Build it with `go build` and install it with `go install`. Get help with
`urleen -h`.

The front-end ships with 2k .woff containing the three icons, a service worker
(so it'll load fast on modern mobile browsers) and a manifest.json so it'll
work like an offline app, even though there's little point in a URL squisher
working offline.

If you don't like the front-end, use a different one, and either let urleen
serve it (with the `-w` flag) or put it behind a server; urleen will respond
to `GET /<id>` with a redirect, and `POST /` by remembering the (JSON)
contents of the request body, provided it's a valid URL. The tiny IDs are
sequential base 62 numbers.

It was built to replace an old Rails 3 app, and to polish my Go chops and to
experiment with some HTML5 features that seem to be gaining browser support.

Feel free to fork and improve, or to criticise with issues.
