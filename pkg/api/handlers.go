package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// writeJSON - вспомогательная функция для отправки JSON ответов
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validateAndFixDate - проверяет и корректирует дату задачи
func validateAndFixDate(task *db.Task) error {
	now := time.Now()
	year, month, day := now.Date()
	now = time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	todayStr := now.Format(dateFormat)
	// Если дата не указана - ставим сегодня
	if task.Date == "" {
		task.Date = todayStr
		return nil
	}

	// Проверяем корректность даты
	date, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %v", err)
	}

	// Если есть правило повторения
	if task.Repeat != "" {
		// Проверяем корректность правила
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		// ВАЖНО: вычисляем следующую дату ТОЛЬКО если дата строго в прошлом
		if date.Before(now) {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
		// Если дата сегодня или в будущем - НЕ МЕНЯЕМ
		return nil
	}

	// Если нет повторения и дата в прошлом - ставим сегодня
	if date.Before(now) {
		task.Date = todayStr
	}

	return nil
}
