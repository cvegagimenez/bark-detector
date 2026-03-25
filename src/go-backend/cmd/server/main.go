package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	mqtt "github.com/cvegagimenez/bark-detector/go-backend/internal/mqtt"
	"github.com/cvegagimenez/bark-detector/go-backend/internal/otel"
)

const (
	defaultBroker   = "tcp://localhost:1883"
	defaultClientID = "go-backend-client"
	defaultTopic    = "bark/metrics"
)

func main() {
	ctx := context.Background()

	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
		if err != nil {
			log.Println(err)
		}
	}()

	broker := envOrDefault("MQTT_BROKER", defaultBroker)
	clientID := envOrDefault("MQTT_CLIENT_ID", defaultClientID)
	topic := envOrDefault("MQTT_TOPIC", defaultTopic)

	log.Printf("Go backend server started. broker=%s topic=%s dt_tenant=%s", broker, topic, os.Getenv("DT_TENANT"))
	mqttClient := mqtt.Connect(broker, clientID)

	if err := mqtt.Subscribe(ctx, mqttClient, topic); err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
