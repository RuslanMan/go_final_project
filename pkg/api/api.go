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

	// Обработчик для /api/tasks (GET) - получаем список задач
	http.HandleFunc("/api/tasks", tasksHandler)

	// Обработчик для /api/task/done (POST) - отметить задачу выполненной
	http.HandleFunc("/api/task/done", doneHandler)
}
