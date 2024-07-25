package main

import (
	"log"

	"github.com/IBM/sarama"
)

var kafkaProducer sarama.SyncProducer
var brokers = []string{"localhost:9092"}

func initKafka() {
	var err error
	// Инициализация продюсера
	kafkaProducer, err = sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		log.Fatalf("Ошибка при создании продюсера: %s", err)
	}
}

func sendMessage(content string, id int32) error {
	msg := &sarama.ProducerMessage{
		Topic: "MsgTopic",
		Value: sarama.StringEncoder(content),
	}
	partition, offset, err := kafkaProducer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Сообщение %s ID: %d отправлено в партицию %d с смещением %d\n", content, id, partition, offset)
	return nil
}
