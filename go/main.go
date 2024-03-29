package main

import (
	"net/http"
)

const (
	webDir = "web"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(":7540", mux)
	if err != nil {
		panic(err)
	}
}
