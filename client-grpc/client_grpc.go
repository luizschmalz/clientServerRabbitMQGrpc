package main

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "cliente-servidor/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewImageServiceClient(conn)

	imageData, err := os.ReadFile("input.png")
	if err != nil {
		panic(err)
	}

	const requests = 50

	// Cria o arquivo de saída
	outputFile, err := os.Create("rtt_grpc.txt")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	for i := 1; i <= requests; i++ {
		start := time.Now()
		resp, err := client.ConvertToGray(context.Background(), &pb.ImageRequest{ImageData: imageData})
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(start)

		// Escreve no arquivo
		_, err = fmt.Fprintf(outputFile, "Requisição %d: %v\n", i, elapsed)
		if err != nil {
			panic(err)
		}

		// Salva a última imagem
		if i == requests {
			err = os.WriteFile("output_grpc.png", resp.ImageData, 0644)
			if err != nil {
				panic(err)
			}
		}
	}
}
