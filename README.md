## Bark Detector Project

### Project Purpose
The goal of this project is to develop a system capable of detecting dog barks using an ESP32 microcontroller. The system processes audio input, identifies bark events, and can be used for monitoring environments where barking is relevant (e.g., kennels, homes, shelters).

### Current Implementation
- Audio acquisition and processing using ESP32
- Bark detection logic
- Communication features (WiFi, SNTP)
- MQTT publishing of measured bark power from the ESP32

### Future Plans
We are extending the project by sending detected bark events and sensor data to an MQTT broker. This allows integration with monitoring and observability platforms, enabling real-time data analysis, alerting, and visualization.

#### Planned Features:
- Integration with external monitoring/observability platforms
- Enhanced data parsing and analytics

### MQTT and Dynatrace flow
- The ESP32 publishes microphone metrics to the MQTT topic `bark/metrics` using the payload format `epoch|sensorID|rms`.
- The Go backend subscribes to the same topic, parses the payload, and exposes the latest `bark_power` metric through OpenTelemetry.
- The local OpenTelemetry collector forwards metrics to Dynatrace tenant `khv5234h`.

To enable Dynatrace export, set `DT_API_TOKEN` in your shell before starting the local stack with Docker Compose.

---
This project is under active development. Contributions and feedback are welcome!
