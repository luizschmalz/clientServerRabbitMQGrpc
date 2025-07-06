package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "cliente-servidor/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewImageServiceClient(conn)

	imageData, err := os.ReadFile("input.png")
	if err != nil {
		log.Fatal(err)
	}

	const requests = 50
	var totalRTT time.Duration

	for i := 1; i <= requests; i++ {
		start := time.Now()
		resp, err := client.ConvertToGray(context.Background(), &pb.ImageRequest{ImageData: imageData})
		if err != nil {
			log.Fatalf("Request %d failed: %v", i, err)
		}
		elapsed := time.Since(start)
		totalRTT += elapsed
		log.Printf("RTT gRPC request %d: %v", i, elapsed)

		// Se quiser salvar só a última resposta, pode fazer aqui:
		if i == requests {
			err = os.WriteFile("output_grpc.png", resp.ImageData, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	avgRTT := totalRTT / requests
	log.Printf("RTT médio gRPC para %d requisições: %v", requests, avgRTT)
}
