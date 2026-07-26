package api

import (
	"net/http"
)

// Init регистрирует все API обработчики
func Init() {
	// Обработчик для /api/nextdate
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// Обработчик для /api/task (POST, GET, PUT, DELETE)
	http.HandleFunc("/api/task", taskHandler)

	// Здесь позже добавим /api/tasks и /api/task/done
	// http.HandleFunc("/api/tasks", tasksHandler)
	// http.HandleFunc("/api/task/done", doneHandler)
}
