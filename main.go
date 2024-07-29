package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// var messge Message
var done = make(chan struct{})

func main() {

	initDB()
	defer db.Close(context.Background())

	initKafka()
	defer kafkaProducer.Close()

	go startConsumer(brokers, os.Getenv("KAFKA_TOPIC"))

	//go launchConsumer(id)

	router := mux.NewRouter()

	router.HandleFunc("/messages", createMessageHandler).Methods("POST")
	router.HandleFunc("/statistics", getStatisticsHandler).Methods("GET")
	// Создаем CORS-миддлвар
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                      // Разрешенные источники, можно указать конкретные домены вместо "*"
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Разрешенные методы
		AllowedHeaders:   []string{"Content-Type"},           // Разрешенные заголовки
		AllowCredentials: true,                               // Разрешить передавать куки
	})

	// Оборачиваем маршрутизатор в CORS-миддлвар
	handler := c.Handler(router)

	// Запускаем сервер с обработчиком CORS
	http.ListenAndServe(":8080", handler)

}

func createMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := saveMessage(message.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message.ID = id
	if err := sendMessage(message.Content, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	<-done
	message.Processed = getCurrentProcessedValue(id)
	writeLine(message, w)
}

func writeLine(message Message, w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(message)
}

// Вывод статистики по обработанным сообщениям
func getStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := getStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}
