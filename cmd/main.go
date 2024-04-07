package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lpernett/godotenv"
	_ "modernc.org/sqlite"
)

const (
	webDir   = "../web/"
	portName = "TODO_PORT"
)

func main() {
	loadEnv()
	connectDB()

	var port string
	port, ok := os.LookupEnv(portName)
	if !ok {
		slog.Info("can't find port in .env.", "port", portName)
		port = ":7540"
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(port, mux)
	if err != nil {
		slog.Error("failed to listen and serve.", "err", err)
		os.Exit(1)
	}
}

func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("can't load .env file.")
		os.Exit(1)
	}
}

func connectDB() {
	appPath, err := os.Executable()
	if err != nil {
		slog.Error("failed to return the path.", "err", err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		slog.Error("failed to connect db.", "err", err)
		os.Exit(1)
	}

	create := `
	    CREATE TABLE scheduler(
			id INTEGER NOT NULL PRIMARY KEY,
			date TEXT(8),
			title TEXT,
			comment TEXT,
			repeat TEXT(128),
		);
		CREATE INDEX date_index ON scheduler (column date);`

	if install {
		if _, err := db.Exec(create); err != nil {
			slog.Error("failed to create db", "err", err)
			os.Exit(1)
		}
	}
}
