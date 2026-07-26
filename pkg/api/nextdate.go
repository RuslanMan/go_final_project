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
// (игнорируем время, сравниваем только даты)
func afterNow(d, now time.Time) bool {
	// Обнуляем время для корректного сравнения
	d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return d.After(now)
}

// NextDate вычисляет следующую дату выполнения задачи
// Поддерживает только базовые правила: d и y
// now - текущая дата
// dstart - начальная дата в формате 20060102
// repeat - правило повторения
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверяем, что правило не пустое
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	// Парсим начальную дату
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата: %v", err)
	}

	// Разбираем правило повторения
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила")
	}

	rule := parts[0]

	switch rule {
	case "d":
		// Правило: d <число> - перенос на указанное количество дней
		// Проверяем формат: должно быть два элемента
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d, требуется: d <число>")
		}

		// Преобразуем интервал в число
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверный интервал для d: %v", err)
		}

		// Проверяем, что интервал в допустимых пределах (1-400)
		if interval <= 0 {
			return "", fmt.Errorf("интервал должен быть положительным числом")
		}
		if interval > 400 {
			return "", fmt.Errorf("интервал не может превышать 400")
		}

		// Сдвигаем дату на interval дней до тех пор,
		// пока она не станет больше now
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		// Правило: y - ежегодное повторение
		// Проверяем формат: должен быть только один элемент
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y, требуется: y")
		}

		// Сдвигаем дату на 1 год до тех пор,
		// пока она не станет больше now
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	default:
		// Если правило не поддерживается, возвращаем ошибку
		return "", fmt.Errorf("неподдерживаемое правило: %s (поддерживаются только d и y)", rule)
	}

	// Возвращаем дату в формате 20060102
	return date.Format(dateFormat), nil
}

// nextDateHandler обрабатывает GET /api/nextdate
// Параметры: now, date, repeat
// Возвращает: просто строку с датой или текст ошибки
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это GET запрос
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	// Если параметр now не указан, используем текущую дату
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			// Возвращаем просто текст ошибки
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("некорректный формат now"))
			return
		}
	}

	// Проверяем наличие обязательных параметров
	if dateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("не указан параметр date"))
		return
	}

	if repeatStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("не указан параметр repeat"))
		return
	}

	// Вызываем функцию NextDate
	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		// Возвращаем просто текст ошибки
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Возвращаем просто строку с датой
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
