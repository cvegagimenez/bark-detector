package client

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"context"
	"log"
	"github.com/cvegagimenez/bark-detector/go-backend/internal/controller"
)

func Connect(broker string, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}
	
	return client
}

func Subscribe(ctx context.Context, client mqtt.Client, topic string) error {
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		err := controller.GetMaxBarkPower(ctx, payload)
		if err != nil {
			log.Printf("Error processing message on topic %s: %v", msg.Topic(), err)
			return
		}
		log.Printf("Received message on topic %s", msg.Topic())
	}

	if token := client.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Subscribed to topic:", topic)

	return nil
}