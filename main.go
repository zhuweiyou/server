package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 3000, "port to listen on")
	flag.Parse()

	http.Handle("/", FileHandler)

	go func() {
		<-time.After(time.Millisecond * 500)
		OpenBrowser(port)
	}()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
