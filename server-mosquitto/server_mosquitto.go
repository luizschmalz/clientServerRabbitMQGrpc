package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"cliente-servidor/imageutils"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("image_processor_server")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	defer client.Disconnect(250)

	if token := client.Subscribe("image/topic", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Mensagem recebida")

		img, err := imageutils.BytesToImage(msg.Payload())
		if err != nil {
			log.Println("Erro decodificando imagem:", err)
			return
		}

		gray := imageutils.ToGray(img)
		grayBytes := imageutils.ImageToBytes(gray)

		token := client.Publish("image/response", 0, false, grayBytes)
		token.Wait()
		if token.Error() != nil {
			log.Println("Erro publicando resposta:", token.Error())
		} else {
			fmt.Println("Imagem processada enviada")
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	select {}
}
