package db

import (
	"database/sql"
	"fmt"
	"go-final-project/internal/types"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Add(task types.Task) (int64, error) {
	slog.Info("Adding new task", "task", task)
	res, err := sq.Insert("scheduler").
		Columns("date", "title", "comment", "repeat").
		Values(task.Date, task.Title, task.Comment, task.Repeat).
		RunWith(r.db).Exec()
	if err != nil {
		return 0, fmt.Errorf("error inserting task into db: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted id: %w", err)
	}
	return int64(id), nil
}

func (r *TaskRepository) GetAllTasks() ([]types.Task, error) {
	slog.Info("Getting all tasks with limit 10")
	rows, err := sq.Select("*").
		From("scheduler").
		OrderBy("date").
		Limit(10).
		RunWith(r.db).Query()
	if err != nil {
		return nil, fmt.Errorf("error getting tasks from db: %w", err)
	}
	defer rows.Close()

	var res []types.Task
	for rows.Next() {
		t := types.Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		res = append(res, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %w", err)
	}
	return res, nil
}

func (r *TaskRepository) GetById(id int64) (types.Task, error) {
	slog.Info("Getting task by id", "id", id)
	row := sq.Select("*").
		From("scheduler").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).QueryRow()
	t := types.Task{}
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	return t, err
}

func (r *TaskRepository) UpdateTask(task types.Task) error {
	slog.Info("Updating task", "task", task)
	res, err := sq.Update("scheduler").
		SetMap(map[string]interface{}{
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}).
		Where(sq.Eq{"id": task.ID}).
		RunWith(r.db).Exec()
	if err != nil {
		return fmt.Errorf("update failed")
	}
	nRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows")
	}
	if nRows == 0 {
		return fmt.Errorf("no rows were updated")
	}
	return nil
}

func (r *TaskRepository) DeleteTask(id int64) error {
	slog.Info("Deleting task by id", "id", id)
	res, err := sq.Delete("scheduler").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).Exec()
	if err != nil {
		return fmt.Errorf("delete failed")
	}
	nRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows")
	}
	if nRows == 0 {
		return fmt.Errorf("no rows were deleted")
	}
	return nil
}
