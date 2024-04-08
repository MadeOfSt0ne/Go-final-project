package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	webDir = "../web/"
)

func StartServer(port string) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("GET /api/nextdate", handleNextDate)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		slog.Error("failed to listen and serve.", "err", err)
		os.Exit(1)
	}
}

func handleNextDate(w http.ResponseWriter, r *http.Request) {
	nowValue := r.FormValue("now")
	dateValue := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if len(repeat) == 0 {
		slog.Debug("empty repeat.")
		http.Error(w, "Empty repeat", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowValue)
	if err != nil {
		slog.Debug("failed to parse time.", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("20060102", dateValue)
	if err != nil {
		slog.Debug("failed to parse time.", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		slog.Debug("failed to get next date.", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}

func NextDate(now, date time.Time, repeat string) (string, error) {
	rule := strings.Split(repeat, " ")
	if len(rule) == 1 && rule[0] != "y" {
		return "", fmt.Errorf("wrong format of repeat: %v", repeat)
	}
	var next time.Time
	switch rule[0] {
	case "d":
		daysToAdd, err := strconv.Atoi(rule[1])
		if err != nil {
			return "", fmt.Errorf("wrong format of repeat: %v", repeat)
		}
		if daysToAdd > 400 {
			return "", fmt.Errorf("max amount of days is 400! Your value is %v", daysToAdd)
		}
		next = date.AddDate(0, 0, daysToAdd)
		for next.Before(now) {
			next = next.AddDate(0, 0, daysToAdd)
		}
	case "y":
		next = date.AddDate(1, 0, 0)
		for next.Before(now) {
			next = next.AddDate(1, 0, 0)
		}
	case "w":

	case "m":

	default:
		return "", fmt.Errorf("wrong format of repeat: %v", repeat)
	}
	return next.Format("20060102"), nil
}
