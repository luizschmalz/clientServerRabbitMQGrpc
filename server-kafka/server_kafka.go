package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"cliente-servidor/imageutils"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "image_topic",
		GroupID: "image_processor_group",
	})

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "image_response",
	})

	defer reader.Close()
	defer writer.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Failed to read message:", err)
		}

		img, err := imageutils.BytesToImage(msg.Value)
		if err != nil {
			log.Println("Invalid image data:", err)
			continue
		}

		gray := imageutils.ToGray(img)
		grayBytes := imageutils.ImageToBytes(gray)

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   msg.Key,
			Value: grayBytes,
		})
		if err != nil {
			log.Println("Failed to write response:", err)
		}
	}
}