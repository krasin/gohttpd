package main

import (
       "flag"
       "fmt"
       "http"
       "io"
       "log"
       "os"
       "path"
       "template"
)

var (
    port = flag.Int("port", 8000, "Port to listen")
    dirTemplate = template.MustParse(`<html>
<head><title>Listing of {Path}</title></head>
<body>
  {.repeated section Names}{@}<br>
  {.end}
</body>
</html>
`, nil)
)

type Dir struct {
     Path string
     Names []string
}

func FileServer(w http.ResponseWriter, req *http.Request) {
     filename := path.Join(".", req.URL.Path)
     log.Printf("Filename: %s", filename)
     f, err := os.Open(filename)
     if err != nil {
     	http.Error(w, "Not found", 404)
	return
     }
     defer f.Close()
     fi, err := f.Stat()
     if err != nil {
     	http.Error(w, "Server error", 500)
	return
     }
     if fi.IsRegular() {
     	io.WriteString(w, fmt.Sprintf("Contents of %s", filename))     	
	return
     }
     if fi.IsDirectory() {
     	names, err := f.Readdirnames(-1)
	if err != nil {
	   http.Error(w, "Server error", 500)
	   return
	}
	err = dirTemplate.Execute(w, &Dir{ filename, names })
	if err != nil {
	   http.Error(w, "Server error", 500)
	   return
	}
	return
     }

     http.Error(w, "Not found", 404)
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
