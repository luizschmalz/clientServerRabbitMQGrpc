package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func main() {
	imageData, err := os.ReadFile("input.png")
	if err != nil {
		log.Fatal("Falha ao ler imagem:", err)
	}

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("image_client_" + uuid.New().String())
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	defer client.Disconnect(250)

	// --- NOVO: abrir arquivo para salvar RTT ---
	outputFile, err := os.Create("mqtt_rtt_resultados.txt")
	if err != nil {
		log.Fatal("Falha ao criar arquivo:", err)
	}
	defer outputFile.Close()

	var wg sync.WaitGroup
	const totalRequests = 50

	pending := make(map[string]chan []byte)
	var mu sync.Mutex

	if token := client.Subscribe("image/response", 0, func(client mqtt.Client, msg mqtt.Message) {
		mu.Lock()
		for id, ch := range pending {
			select {
			case ch <- msg.Payload():
			default:
			}
			delete(pending, id)
			break
		}
		mu.Unlock()
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			corrID := uuid.New().String()
			ch := make(chan []byte)
			mu.Lock()
			pending[corrID] = ch
			mu.Unlock()

			start := time.Now()

			if token := client.Publish("image/topic", 0, false, imageData); token.Wait() && token.Error() != nil {
				log.Fatal("Falha ao enviar mensagem:", token.Error())
			}

			select {
			case resp := <-ch:
				rtt := time.Since(start)
				
				// --- ALTERADO: grava RTT no arquivo em vez de imprimir ---
				mu.Lock()
				_, err := fmt.Fprintf(outputFile, "Requisição %d: RTT %v\n", i+1, rtt)
				mu.Unlock()
				if err != nil {
					log.Println("Erro escrevendo no arquivo:", err)
				}

				if i == totalRequests-1 {
					err := os.WriteFile("output_mqtt.png", resp, 0644)
					if err != nil {
						log.Println("Erro salvando imagem:", err)
					}
				}
			case <-time.After(10 * time.Second):
				mu.Lock()
				_, err := fmt.Fprintf(outputFile, "Requisição %d: timeout\n", i+1)
				mu.Unlock()
				if err != nil {
					log.Println("Erro escrevendo timeout no arquivo:", err)
				}
			}
		}(i)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("Todas as requisições finalizadas.")
}
