package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {
	// Получаем порт из переменной окружения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Получаем путь к файлу БД
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализация базы данных
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	log.Printf("База данных инициализирована: %s", dbFile)

	// ИНИЦИАЛИЗАЦИЯ API ОБРАБОТЧИКОВ
	api.Init()
	log.Println("API обработчики зарегистрированы")

	// Директория с веб-файлами
	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Запуск сервера
	log.Printf("Сервер запущен на http://localhost:%s", port)
	log.Printf("Раздаем файлы из директории: %s", webDir)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
