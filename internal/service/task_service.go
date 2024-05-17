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
	slog.Info("Processing next date with", "now", nowValue, "date", dateValue, "repeat", repeat)
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

	if len(repeat) == 0 {
		return nowValue, nil
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
	//case "w":

	//case "m":

	default:
		return "", fmt.Errorf("wrong format of repeat: %v", repeat)
	}
	return next.Format("20060102"), nil
}

func (s TaskService) AddNewTask(task types.Task) (int64, error) {
	slog.Info("Adding new", "task", task)
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

func (s TaskService) GetTasks() (map[string][]types.TaskDTO, error) {
	slog.Info("Getting tasks")
	tasks, err := s.store.GetAllTasks()
	if err != nil {
		slog.Error("repository returned error.", "err", err)
		return nil, fmt.Errorf("failed to get tasks")
	}
	if tasks == nil {
		tasks = make([]types.Task, 0)
	}

	tasksDTO := make([]types.TaskDTO, 0)
	for _, t := range tasks {
		dto := toTaskDto(t)
		tasksDTO = append(tasksDTO, dto)
	}

	res := make(map[string][]types.TaskDTO)
	res["tasks"] = tasksDTO
	return res, nil
}

func (s TaskService) GetTaskById(taskId string) (types.TaskDTO, error) {
	slog.Info("Getting task by", "id", taskId)
	id, err := strconv.Atoi(taskId)
	if err != nil {
		return types.TaskDTO{}, fmt.Errorf("wrong format of task id: %v", taskId)
	}
	task, err := s.store.GetById(int64(id))
	if err != nil {
		slog.Error("repository returned error.", "err", err)
		return types.TaskDTO{}, fmt.Errorf("failed to get tasks")
	}
	return toTaskDto(task), nil
}

func (s TaskService) UpdateTask(dto types.TaskDTO) error {
	slog.Info("Updating task", "taskDTO", dto)
	if dto.Date != "today" && dto.Date != "" {
		next, err := s.NextDate(time.Now().Format("20060102"), dto.Date, dto.Repeat)
		if err != nil {
			slog.Error("failed to get next date.", "err", err)
			return err
		}
		dto.Date = next
	} else {
		dto.Date = time.Now().Format("20060102")
	}
	if len(dto.Title) == 0 {
		return fmt.Errorf("empty task title")
	}
	if len(dto.ID) == 0 {
		return fmt.Errorf("empty task id")
	}
	id, err := strconv.Atoi(dto.ID)
	if err != nil {
		return fmt.Errorf("wrong format of task id: %v", dto.ID)
	}
	task := types.Task{
		ID:      int64(id),
		Date:    dto.Date,
		Title:   dto.Title,
		Comment: dto.Comment,
		Repeat:  dto.Repeat,
	}
	err = s.store.UpdateTask(task)
	return err
}

func (s TaskService) DeleteTask(taskId string) error {
	slog.Info("Deleting task by", "id", taskId)
	id, err := strconv.Atoi(taskId)
	if err != nil {
		return fmt.Errorf("wrong format of task id: %v", taskId)
	}
	err = s.store.DeleteTask(int64(id))
	if err != nil {
		slog.Error("repository returned error.", "err", err)
		return fmt.Errorf("failed to delete task")
	}
	return nil
}

func (s TaskService) SetNewDate(taskId string) error {
	slog.Info("Setting new date for task", "id", taskId)
	id, err := strconv.Atoi(taskId)
	if err != nil {
		return fmt.Errorf("wrong format of task id: %v", taskId)
	}
	task, err := s.store.GetById(int64(id))
	if err != nil {
		slog.Error("repository returned error.", "err", err)
		return fmt.Errorf("failed to get task")
	}
	if len(task.Repeat) == 0 {
		err = s.store.DeleteTask(int64(id))
		return err
	}

	next, err := s.SetNextDate(task.Date, task.Repeat)
	if err != nil {
		slog.Error("failed to get next date.", "err", err)
		return err
	}
	task.Date = next
	err = s.store.UpdateTask(task)
	if err != nil {
		slog.Error("failed to set next date.", "err", err)
		return err
	}
	return nil
}

func toTaskDto(t types.Task) types.TaskDTO {
	dto := types.TaskDTO{
		ID:      strconv.Itoa(int(t.ID)),
		Title:   t.Title,
		Date:    t.Date,
		Comment: t.Comment,
		Repeat:  t.Repeat,
	}
	return dto
}

func (s TaskService) SetNextDate(dateValue, repeat string) (string, error) {
	date, err := time.Parse("20060102", dateValue)
	if err != nil {
		slog.Error("failed to parse time.", "err", err)
		return "", fmt.Errorf("wrong time format: %v", dateValue)
	}

	rule := strings.Split(repeat, " ")
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
	case "y":
		next = date.AddDate(1, 0, 0)
	default:
		return "", fmt.Errorf("wrong format of repeat: %v", repeat)
	}
	return next.Format("20060102"), nil
}
