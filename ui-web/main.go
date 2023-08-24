package main

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"flag"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/HuguesGuilleus/airazor"
)

const (
	collectionPath = "collection.json"
)

func main() {
	addr := flag.String("a", ":8000", "Listen address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(front)))
	mux.Handle("/api/", &Server{
		Read:         func() ([]byte, error) { return os.ReadFile(collectionPath) },
		Save:         func(data []byte) error { return os.WriteFile(collectionPath, data, 0o664) },
		RoundTripper: http.DefaultTransport,
		LimitBody:    500_000,
		MaxDuration:  time.Millisecond * 500,
	})

	log.Println("listen", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

var (
	//go:embed front/*
	_front embed.FS

	front = func() fs.FS {
		fsys, err := fs.Sub(_front, "front")
		if err != nil {
			panic(err)
		}
		return fsys
	}()
)

type Server struct {
	Read func() ([]byte, error)
	Save func(data []byte) error
	http.RoundTripper
	LimitBody   int64
	MaxDuration time.Duration
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache", "no-store")
	switch r.Method + " " + r.URL.Path {
	case "GET /api/collection.json", "HEAD /api/collection.json":
		data, err := s.Read()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		collection, err := airazor.ParseCollection(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sendCollection(w, collection)

	case "PUT /api/collection.json":
		collection := readCollection(w, r)
		if collection == nil {
			return
		}

		data, _ := json.MarshalIndent(collection, "", "\t")
		s.Save(data)

	case "POST /api/exe":
		collection := readCollection(w, r)
		if collection == nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), s.MaxDuration)
		defer cancel()

		collection.Fetch(&airazor.Config{
			Context:      ctx,
			RoundTripper: s.RoundTripper,
			LimitBody:    s.LimitBody,
		})

		sendCollection(w, collection)
	default:
		http.NotFound(w, r)
	}
}

func readCollection(w http.ResponseWriter, r *http.Request) *airazor.Collection {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	collection, err := airazor.ParseCollection(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return collection
}

func sendCollection(w http.ResponseWriter, collection *airazor.Collection) {
	data, _ := json.Marshal(collection)
	serveBytes(w, "application/json", data)
}

func serveBytes(w http.ResponseWriter, mime string, content []byte) {
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Write(content)
}
