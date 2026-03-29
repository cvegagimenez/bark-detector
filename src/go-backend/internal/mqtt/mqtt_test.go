package client

import (
	"context"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// mockToken implements mqtt.Token with no-op behaviour.
type mockToken struct{}

func (t *mockToken) Wait() bool                       { return true }
func (t *mockToken) WaitTimeout(_ time.Duration) bool { return true }
func (t *mockToken) Done() <-chan struct{}            { ch := make(chan struct{}); close(ch); return ch }
func (t *mockToken) Error() error                     { return nil }

// mockClient implements mqtt.Client, capturing the Subscribe handler for inspection.
type mockClient struct {
	capturedHandler mqtt.MessageHandler
}

func (m *mockClient) IsConnected() bool                                          { return true }
func (m *mockClient) IsConnectionOpen() bool                                     { return true }
func (m *mockClient) Connect() mqtt.Token                                        { return &mockToken{} }
func (m *mockClient) Disconnect(_ uint)                                          {}
func (m *mockClient) Publish(_ string, _ byte, _ bool, _ interface{}) mqtt.Token { return &mockToken{} }
func (m *mockClient) SubscribeMultiple(_ map[string]byte, _ mqtt.MessageHandler) mqtt.Token {
	return &mockToken{}
}
func (m *mockClient) Unsubscribe(_ ...string) mqtt.Token       { return &mockToken{} }
func (m *mockClient) AddRoute(_ string, _ mqtt.MessageHandler) {}
func (m *mockClient) OptionsReader() mqtt.ClientOptionsReader  { return nil }
func (m *mockClient) Subscribe(_ string, _ byte, cb mqtt.MessageHandler) mqtt.Token {
	m.capturedHandler = cb
	return &mockToken{}
}

// mockMessage implements mqtt.Message.
type mockMessage struct {
	topic   string
	payload []byte
}

func (m *mockMessage) Duplicate() bool   { return false }
func (m *mockMessage) Qos() byte         { return 0 }
func (m *mockMessage) Retained() bool    { return false }
func (m *mockMessage) Topic() string     { return m.topic }
func (m *mockMessage) MessageID() uint16 { return 0 }
func (m *mockMessage) Payload() []byte   { return m.payload }
func (m *mockMessage) Ack()              {}

func TestSubscribeReturnsNoError(t *testing.T) {
	client := &mockClient{}
	if err := Subscribe(context.Background(), client, "bark/metrics"); err != nil {
		t.Fatalf("Subscribe() unexpected error: %v", err)
	}
}

func TestSubscribeRegistersHandler(t *testing.T) {
	client := &mockClient{}
	_ = Subscribe(context.Background(), client, "bark/metrics")

	if client.capturedHandler == nil {
		t.Fatal("expected a message handler to be registered, got nil")
	}
}

func TestSubscribeHandlerValidPayload(t *testing.T) {
	client := &mockClient{}
	_ = Subscribe(context.Background(), client, "bark/metrics")

	// A valid payload must not panic or cause any unrecovered error.
	msg := &mockMessage{topic: "bark/metrics", payload: []byte("1711180800|esp32-mic-01|12.34")}
	client.capturedHandler(client, msg)
}

func TestSubscribeHandlerInvalidPayloadDoesNotPanic(t *testing.T) {
	client := &mockClient{}
	_ = Subscribe(context.Background(), client, "bark/metrics")

	// An invalid payload should be logged and discarded gracefully.
	msg := &mockMessage{topic: "bark/metrics", payload: []byte("bad-payload")}
	client.capturedHandler(client, msg)
}
