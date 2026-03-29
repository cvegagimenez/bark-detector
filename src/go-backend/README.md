# Go Backend

Subscribes to bark power measurements from an MQTT broker and exports them as OpenTelemetry metrics.

## Structure

```
├── cmd/server/          # Application entry point
└── internal/
    ├── controller/      # Payload parsing and metric recording
    ├── mqtt/            # MQTT client and subscription handler
    └── otel/            # OpenTelemetry SDK setup and metric gauge
```

## Running locally

```bash
export DT_TENANT=your-tenant-id
export DT_API_TOKEN=dt0c01.XXXX.YYYY
export MQTT_BROKER=tcp://localhost:1883

go run ./cmd/server
```

## Running tests

```bash
go test ./...
```

## Environment variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DT_TENANT` | Yes | — | Dynatrace tenant ID |
| `DT_API_TOKEN` | Yes | — | Dynatrace API token with `metrics.ingest` scope |
| `MQTT_BROKER` | No | `tcp://localhost:1883` | MQTT broker address |
| `MQTT_CLIENT_ID` | No | `go-backend-client` | MQTT client identifier |
| `MQTT_TOPIC` | No | `bark/metrics` | MQTT topic to subscribe to |
