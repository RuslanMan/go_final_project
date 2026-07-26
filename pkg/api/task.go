package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go_final_project/pkg/db"
)

// taskHandler - главный обработчик для /api/task
// Перенаправляет запросы в зависимости от HTTP метода
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
// Создает новую задачу
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// 1. Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "некорректный JSON"})
		return
	}

	// 2. Проверяем обязательное поле title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// 3. Проверяем и корректируем дату
	if err := validateAndFixDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// 4. Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// 5. Возвращаем ID созданной задачи
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

// getTaskHandler - GET /api/task?id=...
// Получает задачу по ID (будет реализован на шаге 6)
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Временная заглушка
	writeJSON(w, map[string]string{"error": "метод GET еще не реализован"})
}

// updateTaskHandler - PUT /api/task
// Обновляет существующую задачу (будет реализован на шаге 6)
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Временная заглушка
	writeJSON(w, map[string]string{"error": "метод PUT еще не реализован"})
}

// deleteTaskHandler - DELETE /api/task?id=...
// Удаляет задачу по ID (будет реализован на шаге 7)
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Временная заглушка
	writeJSON(w, map[string]string{"error": "метод DELETE еще не реализован"})
}
