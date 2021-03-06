package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed dist/*
var dist embed.FS

type distFS struct {
	content embed.FS
}

func (f distFS) Open(name string) (fs.File, error) {
	return f.content.Open(path.Join("dist", name))
}

func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Add("Cache-Control", "public, max-age=31536000, immutable")
		}
		fs.ServeHTTP(w, r)
	}
}

var FileHandler = addHeaders(http.FileServer(http.FS(distFS{dist})))
