package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>math.fent.uk</h1>")
	fmt.Fprintln(w, "<p>ok</p>")
}

func main() {
	cgi.Serve(http.HandlerFunc(handler))
}

