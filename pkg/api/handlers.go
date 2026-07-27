package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// writeJSON - вспомогательная функция для отправки JSON ответов
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

// writeError - вспомогательная функция для отправки ошибок
func writeError(w http.ResponseWriter, message string, statusCode int) {
	writeJSON(w, map[string]string{"error": message}, statusCode)
}

// writeInternalError - возвращает стандартную ошибку для внутренних проблем
func writeInternalError(w http.ResponseWriter) {
	writeError(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
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
		return fmt.Errorf("некорректный формат даты: %w", err)
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
