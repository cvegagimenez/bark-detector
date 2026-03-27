package otel

import (
	"testing"
)

func resetState(t *testing.T) {
	t.Helper()
	stateMu.Lock()
	barkValueByID = map[string]float64{}
	stateMu.Unlock()
	t.Cleanup(func() {
		stateMu.Lock()
		barkValueByID = map[string]float64{}
		stateMu.Unlock()
	})
}

func TestRecordBarkPower(t *testing.T) {
	resetState(t)

	RecordBarkPower(42.5, "sensor-01")

	stateMu.RLock()
	got := barkValueByID["sensor-01"]
	stateMu.RUnlock()

	if got != 42.5 {
		t.Fatalf("expected 42.5, got %f", got)
	}
}

func TestRecordBarkPowerMultipleSensors(t *testing.T) {
	resetState(t)

	RecordBarkPower(10.0, "sensor-01")
	RecordBarkPower(20.0, "sensor-02")

	stateMu.RLock()
	v1 := barkValueByID["sensor-01"]
	v2 := barkValueByID["sensor-02"]
	stateMu.RUnlock()

	if v1 != 10.0 {
		t.Fatalf("sensor-01: expected 10.0, got %f", v1)
	}
	if v2 != 20.0 {
		t.Fatalf("sensor-02: expected 20.0, got %f", v2)
	}
}

func TestRecordBarkPowerOverwritesSameSensor(t *testing.T) {
	resetState(t)

	RecordBarkPower(10.0, "sensor-01")
	RecordBarkPower(99.9, "sensor-01")

	stateMu.RLock()
	got := barkValueByID["sensor-01"]
	stateMu.RUnlock()

	if got != 99.9 {
		t.Fatalf("expected 99.9, got %f", got)
	}
}
