package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	mqtt "github.com/cvegagimenez/bark-detector/go-backend/internal/mqtt"
	"github.com/cvegagimenez/bark-detector/go-backend/internal/otel"
)

const (
	broker   = "tcp://localhost:1883"
	clientID = "go-backend-client"
	topic    = "test"
)

func main() {
	ctx := context.Background()

	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		log.Fatal(err)
	}

    defer func() {
        err = errors.Join(err, otelShutdown(ctx))
        log.Println(err)
    }()

	fmt.Println("Go backend server started")
	mqttClient := mqtt.Connect(broker, clientID)

	if err := mqtt.Subscribe(ctx, mqttClient, topic); err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
