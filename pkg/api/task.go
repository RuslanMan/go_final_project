package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go_final_project/pkg/db"
)

// taskHandler - главный обработчик для /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// addTaskHandler - POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "некорректный JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := validateAndFixDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

// getTaskHandler - GET /api/task?id=...
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, task)
}

// updateTaskHandler - PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "некорректный JSON"})
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := validateAndFixDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}

// deleteTaskHandler - DELETE /api/task?id=...
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	// Удаляем задачу из БД
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]string{})
}
