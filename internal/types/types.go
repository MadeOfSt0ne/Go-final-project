package types

// Task structure
type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Task DTO structure
type TaskDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Task store interface
type TaskStore interface {
	Add(task Task) (int64, error)
	GetAllTasks() ([]Task, error)
	GetById(id int64) (Task, error)
	UpdateTask(task Task) error
	DeleteTask(id int64) error
}

// Response structure
type ResponseOK struct {
	ID string `json:"id"`
}
