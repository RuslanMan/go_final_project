package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Константа для формата даты
const dateFormat = "20060102"

// afterNow проверяет, что дата d больше даты now
func afterNow(d, now time.Time) bool {
	d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return d.After(now)
}

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата: %w", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила")
	}

	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d, требуется: d <число>")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверный интервал для d: %w", err)
		}

		if interval <= 0 {
			return "", fmt.Errorf("интервал должен быть положительным числом")
		}
		if interval > 400 {
			return "", fmt.Errorf("интервал не может превышать 400")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y, требуется: y")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("неподдерживаемое правило: %s (поддерживаются только d и y)", rule)
	}

	return date.Format(dateFormat), nil
}

// nextDateHandler обрабатывает GET /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "некорректный формат now", http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, "не указан параметр date", http.StatusBadRequest)
		return
	}

	if repeatStr == "" {
		http.Error(w, "не указан параметр repeat", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
