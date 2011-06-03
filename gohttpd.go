package main

import (
	"flag"
	"fmt"
	"http"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"template"
)

var (
	port        = flag.Int("port", 8000, "Port to listen")
	dirTemplate = template.MustParse(`<html>
<head><title>Listing of {Path}</title></head>
<body>
  {.repeated section Names}<a href="/{Path}/{@}">{@}</a><br>
  {.end}
</body>
</html>
`, nil)

	mimeTypes = map[string]string{
		"gif":  "image/gif",
		"jpeg": "image/jpeg",
		"jpg":  "image/jpg",
		"png":  "image/png",
		"svg":  "image/svg+xml",
		"tiff": "image/tiff",
		"mp3":  "audio/mp3",
		"pdf":  "application/pdf",
		"js":   "application/javascript",
		"zip":  "application/zip",
		"gz":   "application/x-gzip",
		"json": "application/json",
		"nexe": "application/x-nacl",
		"txt":  "text/plain",
		"html": "text/html",
		"htm":  "text/html",
		"css":  "text/css",
		"csv":  "text/csv",
		"xml":  "text/xml",
		"tar":  "application/x-tar",
		"go":   "text/plain",
	}

	knownFiles = map[string]string{
		"Makefile":   "text/plain",
		".gitignore": "text/plain",
		"README":     "text/plain",
		"LICENSE":    "text/plain",
		"configure":  "text/plain",
	}
)

type Dir struct {
	Path  string
	Names []string
}

func guessMimeType(name string) string {
	ext := path.Ext(strings.ToLower(name))
	if ext != "" {
		ext = ext[1:]
	}
	typ, ok := mimeTypes[ext]
	if ok {
		return typ
	}
	typ, ok = knownFiles[name]
	if ok {
		return typ
	}
	return "application/octet-stream"
}

func sendFile(w http.ResponseWriter, f *os.File) {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	writeHeaders(w, guessMimeType(f.Name()), len(data))
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
}

func writeHeaders(w http.ResponseWriter, mimeType string, length int) {
     h := w.Header()
     h.Set("Cache-Control", "no-cache")    		 
     h.Set("Content-Type", mimeType)
     h.Set("Content-Length", strconv.Itoa(length))
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
		sendFile(w, f)
		return
	}
	if fi.IsDirectory() {
		names, err := f.Readdirnames(-1)
		if err != nil {
			http.Error(w, "Server error", 500)
			return
		}
		err = dirTemplate.Execute(w, &Dir{filename, names})
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
