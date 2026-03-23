package controller

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cvegagimenez/bark-detector/go-backend/internal/otel"
)

type Measurement struct {
	Timestamp time.Time
	SensorID  string
	BarkPower float64
}

func ParseMetricPayload(payload string) (Measurement, error) {
	splittedPayload := strings.Split(payload, "|")
	if len(splittedPayload) != 3 {
		return Measurement{}, fmt.Errorf("invalid payload format %q: expected epoch|sensorID|barkPower", payload)
	}

	epochInt, err := strconv.ParseInt(strings.TrimSpace(splittedPayload[0]), 10, 64)
	if err != nil {
		return Measurement{}, fmt.Errorf("convert epoch time: %w", err)
	}

	sensorID := strings.TrimSpace(splittedPayload[1])
	if sensorID == "" {
		return Measurement{}, fmt.Errorf("sensor ID is empty")
	}

	barkPower, err := strconv.ParseFloat(strings.TrimSpace(splittedPayload[2]), 64)
	if err != nil {
		return Measurement{}, fmt.Errorf("convert bark power: %w", err)
	}

	return Measurement{
		Timestamp: time.Unix(epochInt, 0),
		SensorID:  sensorID,
		BarkPower: barkPower,
	}, nil
}

func RecordMeasurement(measurement Measurement) {
	log.Printf("Message received at %s: bark_power=%f sensor_id=%s", measurement.Timestamp.Format(time.RFC3339), measurement.BarkPower, measurement.SensorID)

	otel.RecordBarkPower(measurement.BarkPower, measurement.SensorID)
}
