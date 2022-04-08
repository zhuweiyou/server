package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed dist/*
var root embed.FS

type myFS struct {
	content embed.FS
}

func (c myFS) Open(name string) (fs.File, error) {
	return c.content.Open(path.Join("dist", name))
}

func main() {
	http.Handle("/", addHeaders(http.FileServer(http.FS(myFS{root}))))
	http.ListenAndServe(":3000", nil)
}

func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") || strings.HasPrefix(r.URL.Path, "/favicon.ico") {
			w.Header().Add("Cache-Control", "max-age=88888888")
		}
		fs.ServeHTTP(w, r)
	}
}
