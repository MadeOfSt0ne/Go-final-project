package types

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskStore interface {
	Add(task Task) (int64, error)
}

type ResponseOK struct {
	ID string `json:"id"`
}
