package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// Глобальная переменная для доступа к БД из других пакетов
var DB *sql.DB

// Схема базы данных: создание таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init инициализирует подключение к базе данных
// Если файла БД не существует, создает его с таблицей и индексом
func Init(dbFile string) error {
	// Проверяем существование файла БД
	_, err := os.Stat(dbFile)
	needCreate := os.IsNotExist(err)

	// Открываем базу данных
	// Если файл не существует, sqlite создаст его автоматически
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Сохраняем соединение в глобальной переменной
	DB = db

	// Проверяем, что соединение работает
	if err := db.Ping(); err != nil {
		return err
	}

	// Если файла не существовало, создаем таблицу и индекс
	if needCreate {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
