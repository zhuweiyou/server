package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

//go:embed dist/*
var root embed.FS

type staticFS struct {
	content embed.FS
}

func (s staticFS) Open(name string) (fs.File, error) {
	return s.content.Open(path.Join("dist", name))
}

func main() {
	port := flag.Int("port", 3000, "port to listen on")
	flag.Parse()

	http.Handle("/", addHeaders(http.FileServer(http.FS(staticFS{root}))))

	go func() {
		url := fmt.Sprintf("http://localhost:%d/", *port)
		log.Printf("Open %s in your browser", url)

		<-time.After(time.Second)
		_ = openBrowser(url)

	}()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") || strings.HasPrefix(r.URL.Path, "/favicon.ico") {
			w.Header().Add("Cache-Control", "max-age=88888888")
		}
		fs.ServeHTTP(w, r)
	}
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
