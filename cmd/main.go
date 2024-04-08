package main

import (
	"database/sql"
	"go-final-project/internal/api"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lpernett/godotenv"
	_ "modernc.org/sqlite"
)

const (
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
	api.StartServer(port)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("failed to load .env file.", "err", err)
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

	db, err := sql.Open("sqlite", "../scheduler.db")
	if err != nil {
		slog.Error("failed to connect db.", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	create := `
	    CREATE TABLE scheduler(
			id INTEGER PRIMARY KEY,
			date VARCHAR(8),
			title VARCHAR,
			comment VARCHAR,
			repeat VARCHAR(128)
		);
		CREATE INDEX date_idx ON scheduler (date);
		`

	if !install {
		if _, err := db.Exec(create); err != nil {
			slog.Error("failed to create db.", "err", err)
			os.Exit(1)
		}
	}
}
