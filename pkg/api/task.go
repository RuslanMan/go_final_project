package api

import (
	"encoding/json"
	"fmt"
	"log"
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
		writeError(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := validateAndFixDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Ошибка добавления задачи: %v", err)
		writeInternalError(w)
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)}, http.StatusOK)
}

// getTaskHandler - GET /api/task?id=...
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "задача не найдена", http.StatusNotFound)
		return
	}

	writeJSON(w, task, http.StatusOK)
}

// updateTaskHandler - PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := validateAndFixDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		log.Printf("Ошибка обновления задачи (ID=%s): %v", task.ID, err)
		writeInternalError(w)
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}

// deleteTaskHandler - DELETE /api/task?id=...
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Удаляем задачу из БД
	if err := db.DeleteTask(id); err != nil {
		log.Printf("Ошибка удаления задачи (ID=%s): %v", id, err)
		writeInternalError(w)
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
