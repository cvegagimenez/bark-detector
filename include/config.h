#ifndef CONFIG_H
#define CONFIG_H

const long serialBaudRate = 115200;
const int micPin = 35;
const int numberMeasures = 100;

// WiFi config
char wifi_ssid[] = "MIWIFI_mnQN";
char wifi_password[] = "jrcCHfac";
const char mqtt_broker[] = "192.168.1.227";
const int mqtt_port = 1883;
const char mqtt_client_id[] = "bark-detector-micro";
const char mqtt_topic[] = "bark/metrics";
const char sensor_id[] = "esp32-mic-01";

// SNTP config
const int NTP_PACKET_SIZE = 48;
const int udpPort = 2390;
const int ntpPort = 123;
const int daylightOffset_sec = 3600;

// Function prototypes
void beginUdp();
void connectWifi();
void connectMqtt();
void ensureMqttConnected();
bool publishMetric(long epochTime, float rms);
void readMeasure();
long fetchTimeNtp();

#endif
