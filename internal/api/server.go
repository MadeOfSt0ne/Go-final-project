package api

import (
	"database/sql"
	"go-final-project/internal/db"
	"go-final-project/internal/service"
	"log/slog"
	"net/http"
	"os"
)

const (
	webDir = "./web/"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() {
	slog.Info("Running server")
	mux := http.NewServeMux()

	taskStore := db.NewTaskRepository(s.db)
	taskService := service.NewTaskService(taskStore)
	taskHandler := NewHandler(taskService)
	taskHandler.RegisterRoutes(mux)

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(s.addr, mux)
	if err != nil {
		slog.Error("failed to listen and serve.", "err", err)
		os.Exit(1)
	}
}
