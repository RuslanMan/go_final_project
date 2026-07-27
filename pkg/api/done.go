package api

import (
	"log"
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
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "задача не найдена", http.StatusNotFound)
		return
	}

	// Если нет повторения - удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			log.Printf("Ошибка удаления задачи при выполнении (ID=%s): %v", id, err)
			writeInternalError(w)
			return
		}
		writeJSON(w, map[string]string{}, http.StatusOK)
		return
	}

	// Если есть повторение - вычисляем следующую дату
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем дату задачи
	if err := db.UpdateDate(id, next); err != nil {
		log.Printf("Ошибка обновления даты задачи при выполнении (ID=%s): %v", id, err)
		writeInternalError(w)
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
