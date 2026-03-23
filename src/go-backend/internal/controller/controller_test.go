package controller

import (
	"testing"
	"time"
)

func TestParseMetricPayload(t *testing.T) {
	t.Parallel()

	measurement, err := ParseMetricPayload("1711180800|esp32-mic-01|12.34")
	if err != nil {
		t.Fatalf("ParseMetricPayload() error = %v", err)
	}

	if measurement.SensorID != "esp32-mic-01" {
		t.Fatalf("unexpected sensor ID: %s", measurement.SensorID)
	}

	if measurement.BarkPower != 12.34 {
		t.Fatalf("unexpected bark power: %f", measurement.BarkPower)
	}

	expectedTime := time.Unix(1711180800, 0)
	if !measurement.Timestamp.Equal(expectedTime) {
		t.Fatalf("unexpected timestamp: %s", measurement.Timestamp)
	}
}

func TestParseMetricPayloadRejectsInvalidFormat(t *testing.T) {
	t.Parallel()

	if _, err := ParseMetricPayload("1711180800|12.34"); err == nil {
		t.Fatal("expected invalid payload format error")
	}
}

func TestParseMetricPayloadRejectsEmptySensorID(t *testing.T) {
	t.Parallel()

	if _, err := ParseMetricPayload("1711180800||12.34"); err == nil {
		t.Fatal("expected empty sensor ID error")
	}
}
