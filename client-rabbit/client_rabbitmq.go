package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErr(err, "Failed to connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErr(err, "Failed to open channel")
	defer ch.Close()

	replyQueue, err := ch.QueueDeclare("", false, true, true, false, nil)
	failOnErr(err, "Failed to declare reply queue")

	msgs, err := ch.Consume(replyQueue.Name, "", true, false, false, false, nil)
	failOnErr(err, "Failed to consume reply queue")

	imageData, err := os.ReadFile("input.png")
	failOnErr(err, "Failed to read input image")

	// Cria/abre o arquivo de saída
	outputFile, err := os.Create("rtt_resultados.txt")
	failOnErr(err, "Failed to create output file")
	defer outputFile.Close()

	const totalRequests = 50

	for i := 1; i <= totalRequests; i++ {
	corrID := uuid.New().String()
	start := time.Now()

	err = ch.Publish("", "image_queue", false, false, amqp.Publishing{
		ContentType:   "application/octet-stream",
		CorrelationId: corrID,
		ReplyTo:       replyQueue.Name,
		Body:          imageData,
	})
	failOnErr(err, "Failed to publish request")

	for d := range msgs {
		if d.CorrelationId == corrID {
			elapsed := time.Since(start)
			_, err := fmt.Fprintf(outputFile, "Requisição %d: %v\n", i, elapsed)
			failOnErr(err, "Failed to write to file")

			// Salva a imagem apenas na última iteração
			if i == totalRequests {
				err := os.WriteFile("output_rabbitmq.png", d.Body, 0644)
				failOnErr(err, "Failed to save output image")
			}
			break
		}
	}
}}