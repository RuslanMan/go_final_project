package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Получаем порт из переменной окружения TODO_PORT
	// Если переменная не задана, используем порт 7540 по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Директория с веб-файлами
	webDir := "./web"

	// Создаем файловый сервер для раздачи статики
	// http.FileServer будет автоматически обрабатывать запросы к:
	// - / → index.html
	// - /css/style.css → ./web/css/style.css
	// - /js/scripts.min.js → ./web/js/scripts.min.js
	// - /favicon.ico → ./web/favicon.ico
	fileServer := http.FileServer(http.Dir(webDir))

	// Регистрируем обработчик для корневого пути "/"
	// Все запросы будут направляться к файловому серверу
	http.Handle("/", fileServer)

	// Запускаем HTTP сервер
	log.Printf("Сервер запущен на http://localhost:%s", port)
	log.Printf("Раздаем файлы из директории: %s", webDir)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
