package controller

import (
	"log"
	"strconv"
	"strings"
	"time"
	"context"

	"github.com/cvegagimenez/bark-detector/go-backend/internal/otel"
)

var (
	maxBarkPower float64
)

func GetMaxBarkPower(ctx context.Context, payload string) error {
	splittedPayload := strings.Split(payload, "|")

	epochInt, err := strconv.ParseInt(splittedPayload[0], 10, 64)
	if err != nil {
		log.Fatalf("Error converting epoch time: %v", err)
		return err
	}

	sensorID := splittedPayload[1]
	if sensorID == "" {
		log.Fatalf("Error converting sensor ID: %v", err)
		return err
	}

	barkPower, err := strconv.ParseFloat(splittedPayload[2], 64)
	if err != nil {
		log.Fatalf("Error converting bark power: %v", err)
		return err
	}

    if barkPower > maxBarkPower {
        maxBarkPower = barkPower
    }

    parsedTime := time.Unix(epochInt, 0)

    log.Printf("Message received at %s: %v", parsedTime, maxBarkPower)

    otel.RecordBarkPower(maxBarkPower, sensorID)

    return nil
}
