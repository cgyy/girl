package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	Id        int
	UserId    int
	Title     string
	Content   string
	Active    bool
	CreatedAt string
	UpdatedAt string
}

func FindTask(id int) (t Task) {
	row := db.QueryRow("SELECT id, user_id, Title, content, created_at, updated_at FROM tasks WHERE id = ?", id)

	row.Scan(&t.Id, &t.UserId, &t.Title, &t.Content, &t.CreatedAt, &t.UpdatedAt)

	return
}

func FindTasks() (tasks []Task) {
	rows, err := db.Query("SELECT id, user_id, Title, content, created_at, updated_at FROM tasks ORDER BY id ASC limit 100")
	checkErr(err)

	for rows.Next() {
		t := Task{}
		rows.Scan(&t.Id, &t.UserId, &t.Title, &t.Content, &t.CreatedAt, &t.UpdatedAt)
		tasks = append(tasks, t)
	}
	return
}

func DeleteTask(id int) {
	_, err := db.Exec("DELETE FROM tasks WHERE id=?", id)
	fmt.Println("delete", id)
	checkErr(err)
}

// insert
// update
func (r *Task) Save() *Task {
	now := FormatTime(time.Now())
	var err error
	r.UpdatedAt = now
	if r.Id > 0 {
		_, err = db.Exec("UPDATE tasks SET user_id=?, title=?, content=?, updated_at=?  WHERE id=?",
			r.UserId, r.Title, r.Content, now, r.Id)

	} else {
		var insertId int64
		var result sql.Result
		r.CreatedAt = now
		result, err = db.Exec("INSERT INTO tasks(user_id, title, content, created_at, updated_at) values (?,?,?,?,?)",
			r.UserId, r.Title, r.Content, now, now)

		insertId, err = result.LastInsertId()
		r.Id = int(insertId)
	}

	checkErr(err)
	return r
}

func (r *Task) Delete() {
}
