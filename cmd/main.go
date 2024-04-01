package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/lpernett/godotenv"
)

const (
	webDir   = "../web/"
	portName = "TODO_PORT"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("can't load .env file")
		os.Exit(1)
	}

	var port string
	port, ok := os.LookupEnv(portName)
	if !ok {
		slog.Info("can't find port in .env", "port", portName)
		port = ":7540"
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(port, mux)
	if err != nil {
		slog.Error("failed to listen and serve", "error", err)
		os.Exit(1)
	}
}
