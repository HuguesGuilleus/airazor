package main

import (
	_ "embed"
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {
	addr := flag.String("a", ":8000", "Listen address")
	flag.Parse()

	s := Server{}

	log.Println("listen", *addr)
	log.Fatal(http.ListenAndServe(*addr, &s))
}

var (
	//go:embed front/index.html
	indexHTML []byte

	//go:embed front/style.css
	styleCSS []byte

	//go:embed front/index.js
	indexJs []byte
)

type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		serveBytes(w, "text/html", indexHTML)
	case "/index.js":
		serveBytes(w, "application/javascript", indexJs)
	case "/style.css":
		serveBytes(w, "text/css", styleCSS)
	default:
		http.NotFound(w, r)
	}
}

func serveBytes(w http.ResponseWriter, mime string, content []byte) {
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Write(content)
}
