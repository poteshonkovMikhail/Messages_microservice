package main

import (
	"context"

	"log"

	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func initDB() {
	var err error
	db, err = pgx.Connect(context.Background(), "postgres://postgres:workout+5@localhost:5432/messages_microservice_db")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	sql := `CREATE TABLE IF NOT EXISTS messages (
        id SERIAL PRIMARY KEY,
        content TEXT NOT NULL,
        processed BOOLEAN DEFAULT FALSE
    );`

	_, err = db.Exec(context.Background(), sql)
	if err != nil {
		log.Fatalf("unable to create table: %v", err)
	}
}

func saveMessage(content string) (int32, error) {
	var id int32
	err := db.QueryRow(context.Background(), "INSERT INTO messages(content) VALUES($1) RETURNING id", content).Scan(&id)
	setID(id)
	return id, err
}

func updateMessage(id int32, done chan struct{}) error {
	_, err := db.Exec(context.Background(), "UPDATE messages SET processed = TRUE WHERE id = $1", id)
	done <- struct{}{}

	return err
}

func getCurrentProcessedValue(id int32) bool {
	var value bool = false
	sql_statement := "SELECT processed from messages WHERE id = $1"
	err := db.QueryRow(context.Background(), sql_statement, id).Scan(&value)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Fatalf("Ошибка при попытки получить сообщение об успешной обработке: Статус отправленного сообщения не найден (%s)", err) // Если строки нет, возвращаем пустую строку
		}
		log.Fatalf("Ошибка при попытки получить сообщение об успешной обработке: %s", err)
	}

	return value

}

// Получение статистики по обработанным сообщениям
func getStatistics() (MessageStatistics, error) {
	var stats MessageStatistics
	err := db.QueryRow(context.Background(),
		"SELECT COUNT(*), SUM(CASE WHEN processed THEN 1 ELSE 0 END) FROM messages").
		Scan(&stats.TotalMessages, &stats.ProcessedMessages)

	if err != nil {
		return stats, err
	}

	stats.UnprocessedMessages = stats.TotalMessages - stats.ProcessedMessages
	return stats, nil
}
