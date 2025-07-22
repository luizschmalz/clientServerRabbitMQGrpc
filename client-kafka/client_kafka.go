package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	imageData, err := os.ReadFile("input.png")
	failOnErr(err, "Failed to read image")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "image_topic",
	})
	defer writer.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "image_response",
		GroupID: uuid.New().String(), // única para evitar consumo duplicado
	})
	defer reader.Close()

	outputFile, err := os.Create("kafka_rtt_resultados.txt")
	failOnErr(err, "Failed to create RTT file")
	defer outputFile.Close()

	const totalRequests = 50

	for i := 1; i <= totalRequests; i++ {
		corrID := uuid.New().String()
		start := time.Now()

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(corrID),
			Value: imageData,
		})
		failOnErr(err, "Failed to send Kafka message")

		for {
			msg, err := reader.ReadMessage(context.Background())
			failOnErr(err, "Failed to read Kafka response")

			if string(msg.Key) == corrID {
				rtt := time.Since(start)
				fmt.Fprintf(outputFile, "Requisição %d: %v\n", i, rtt)

				if i == totalRequests {
					err := os.WriteFile("output_kafka.png", msg.Value, 0644)
					failOnErr(err, "Failed to save final image")
				}
				break
			}
		}
	}
}