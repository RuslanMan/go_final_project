package api

import (
	"net/http"
)

// Init регистрирует все API обработчики
func Init() {
	// Регистрируем обработчик для /api/nextdate
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// Здесь позже будут добавлены другие обработчики:
	// http.HandleFunc("/api/task", taskHandler)
	// http.HandleFunc("/api/tasks", tasksHandler)
}
