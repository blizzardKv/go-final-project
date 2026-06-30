package db

type Task struct {
	ID      int64  `json:"id"      db:"id"`
	Date    string `json:"date"    db:"date"`
	Title   string `json:"title"   db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat"  db:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
