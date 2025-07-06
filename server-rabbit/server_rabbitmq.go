package main

import (
	"log"
	"cliente-servidor/imageutils"
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

	q, _ := ch.QueueDeclare("image_queue", false, false, false, false, nil)

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			img, _ := imageutils.BytesToImage(d.Body)
			gray := imageutils.ToGray(img)
			grayBytes := imageutils.ImageToBytes(gray)

			ch.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				CorrelationId: d.CorrelationId,
				Body:          grayBytes,
			})
		}
	}()

	log.Println("Waiting for messages (RabbitMQ)...")
	<-forever
}
