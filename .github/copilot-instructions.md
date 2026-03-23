# Copilot Instructions

## Project overview

This repository contains the **Bark Detector** project.

- The main target is an **ESP32** running **Arduino via PlatformIO**.
- The firmware reads microphone samples, computes audio power, gets time from **SNTP**, and is being evolved toward bark-event reporting.
- The repo also includes a small **Go backend** that subscribes to **MQTT** messages and records bark power with **OpenTelemetry**.
- Local supporting services are defined with Docker Compose for **Mosquitto MQTT** and an **OTel collector**.

## Repository structure

- `README.md`: project summary and roadmap.
- `platformio.ini`: PlatformIO configuration for `esp32doit-devkit-v1`.
- `include/config.h`: shared firmware configuration and function declarations. Treat it as sensitive because it may contain local/dev credentials or machine-specific settings.
- `src/micro/`: ESP32 firmware split by concern:
  - `04bark.ino`: Arduino `setup()` / `loop()` entrypoint.
  - `01communication.ino`: WiFi and basic connectivity setup.
  - `02sntp.ino`: NTP/SNTP time sync.
  - `03audio.ino`: microphone sampling and RMS calculation.
- `src/go-backend/`: Go service:
  - `cmd/server/main.go`: backend entrypoint.
  - `internal/mqtt/`: MQTT connection and subscription logic.
  - `internal/controller/`: payload parsing and bark-power handling.
  - `internal/otel/`: telemetry setup/recording.
- `env/`: local config mounted into supporting containers:
  - `env/mqtt/`
  - `env/otel/`
- `docker-compose.yaml`: local MQTT and telemetry stack.
- `test/`: PlatformIO test directory placeholder.

## Working guidance

- Keep the current separation of concerns between the numbered `.ino` files.
- Prefer small, targeted changes that match existing Arduino/PlatformIO and Go patterns.
- Reuse `include/config.h` for shared firmware declarations instead of duplicating constants or prototypes.
- Avoid committing real credentials, hostnames, or other environment-specific values when editing config files.
