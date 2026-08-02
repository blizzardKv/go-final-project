package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      int64  `json:"id,string" db:"id"`
	Date    string `json:"date"       db:"date"`
	Title   string `json:"title"      db:"title"`
	Comment string `json:"comment"    db:"comment"`
	Repeat  string `json:"repeat"     db:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(search string, limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)
	var rows *sql.Rows
	var err error

	switch {
	case search == "":
		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`, limit)
	case isDateSearch(search):
		t, _ := time.Parse("02.01.2006", search)
		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`,
			t.Format("20060102"), limit)
	default:
		like := "%" + search + "%"
		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
			like, like, limit)
	}
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err = rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return tasks, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

func isDateSearch(s string) bool {
	_, err := time.Parse("02.01.2006", s)
	return err == nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
