package db

import (
	"database/sql"
	"fmt"
	"go-final-project/internal/types"

	sq "github.com/Masterminds/squirrel"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Add(task types.Task) (int64, error) {
	res, err := sq.Insert("scheduler").
		Columns("date", "title", "comment", "repeat").
		Values(task.Date, task.Title, task.Comment, task.Repeat).
		RunWith(r.db).Exec()
	if err != nil {
		return 0, fmt.Errorf("error inserting task into db: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted id: %v", err)
	}
	return int64(id), nil
}
