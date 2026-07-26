package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// doneHandler обрабатывает POST /api/task/done
// Отмечает задачу выполненной:
// - Если нет повторения (repeat пустой) - удаляет задачу
// - Если есть повторение - вычисляет следующую дату и обновляет
func doneHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из параметров запроса
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Если нет правила повторения - удаляем задачу
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	// Если есть правило повторения - вычисляем следующую дату
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Обновляем дату задачи
	if err := db.UpdateDate(id, next); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем успех
	writeJSON(w, map[string]string{})
}
