package db

import (
	"database/sql"
	"fmt"
)

// Task - структура задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в базу данных
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) 
	          VALUES (?, ?, ?, ?)`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Tasks возвращает список задач с ограничением по количеству
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          ORDER BY date 
	          LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE id = ?`

	task := &Task{}
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("задача не найдена")
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTask обновляет задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler 
	          SET date = ?, title = ?, comment = ?, repeat = ? 
	          WHERE id = ?`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(id string, date string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	result, err := DB.Exec(query, date, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
