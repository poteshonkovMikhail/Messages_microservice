package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

var id int32 = 51
var idH int32

func setID(id0 int32) {
	id = id0
}

func getID() int32 {
	return id
}

// Функция обработки сообщения
func processMessage(msg *sarama.ConsumerMessage, topic string, partition int32) {
	idH = getID()
	err := updateMessage(idH, done)
	if err != nil {
		log.Fatalf("Ошибка обработки сообщения: %s", err)
	}
	log.Printf("Обработано сообщение: %s | ID: %d | Тема: %s | Партиция: %d", string(msg.Value), idH, topic, partition)
}

// Функция для создания и управления консюмером
func startConsumer(brokers []string, topic string) {
	// Создаем новый конфиг для консюмера
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Создаем консюмера
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Ошибка создания консюмера: %s", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Ошибка закрытия консюмера: %s", err)
		}
	}()

	// Получаем разделы для указанного топика
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Ошибка получения разделов: %s", err)
	}

	// Обрабатываем каждую партицию
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Не удалось создать консюмера для партиции %d: %s", partition, err)
		}
		defer pc.Close()

		// Запускаем отдельную горутину для обработки сообщений из этой партиции
		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					processMessage(msg, topic, partition)
				case err := <-pc.Errors():
					log.Println(err)
				}
			}
		}(pc)
	}

	// Главный контекст для контролирования завершения программы
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка системных сигналов для корректного завершения
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Ждем сигнала для завершения работы
	select {
	case <-sigs:
		log.Println("Получен сигнал завершения работы")
		cancel()
		close(done)
		log.Println("Программа завершена")
	case <-ctx.Done():
		log.Println("Контекст завершён")
	}
}
