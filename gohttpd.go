package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
)

var (
	port = flag.Int("port", 8000, "Port to listen")
)

func DetectContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch path.Ext(r.URL.Path) {
		case ".json":
			w.Header().Add("Content-Type", "application/json")
		case ".js":
			w.Header().Add("Content-Type", "application/javascript")
		case ".css":
			w.Header().Add("Content-Type", "text/css")
		case ".html":
			w.Header().Add("Content-Type", "text/html")
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	http.Handle("/", DetectContentType(http.FileServer(http.Dir("."))))
	fmt.Printf("Listening tcp %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
