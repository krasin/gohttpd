package main

import (
       "flag"
       "fmt"
       "http"
       "io"
       "log"
       "os"
       "path"
)

var port = flag.Int("port", 8000, "Port to listen")

func FileServer(w http.ResponseWriter, req *http.Request) {
     filename := path.Join(".", req.URL.Path)
     log.Printf("Filename: %s", filename)
     fi, err := os.Stat(filename)
     if err != nil {
     	http.Error(w, "", 404)
	return
     }
     if fi.IsRegular() {
     	
     }

     io.WriteString(w, "hello, world!\n")
}

func main() {
     flag.Parse()
     http.HandleFunc("/", FileServer)
     fmt.Printf("Listening tcp %d\n", *port)
     err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
     if err != nil {
     	log.Fatal("ListenAndServe: ", err.String())
     }
}
