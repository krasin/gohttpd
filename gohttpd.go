package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	port = flag.Int("port", 8000, "Port to listen")
)

func main() {
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Printf("Listening tcp %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
