package service

import (
	"fmt"
	"go-final-project/internal/types"
	"log/slog"

	"strconv"
	"strings"
	"time"
)

type TaskService struct {
	store types.TaskStore
}

func NewTaskService(store types.TaskStore) TaskService {
	return TaskService{store: store}
}

func (s TaskService) NextDate(nowValue, dateValue, repeat string) (string, error) {
	now, err := time.Parse("20060102", nowValue)
	if err != nil {
		slog.Error("failed to parse time.", "err", err)
		return "", fmt.Errorf("wrong time format: %v", nowValue)
	}

	date, err := time.Parse("20060102", dateValue)
	if err != nil {
		slog.Error("failed to parse time.", "err", err)
		return "", fmt.Errorf("wrong time format: %v", dateValue)
	}

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

func (s TaskService) AddNewTask(task types.Task) (int64, error) {
	if len(task.Title) == 0 {
		return 0, fmt.Errorf("empty task title")
	}
	id, err := s.store.Add(task)
	if err != nil {
		slog.Error("repository returned error.", "err", err)
		return 0, fmt.Errorf("failed to add new task")
	}
	return id, nil
}
