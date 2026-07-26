package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

// TasksResp - структура ответа со списком задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET /api/tasks
// Возвращает список задач, отсортированных по дате
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем задачи из БД (максимум 50)
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Если задач нет, создаем пустой слайс, чтобы вернуть {"tasks":[]}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	// Возвращаем список
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
