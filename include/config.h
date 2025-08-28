#ifndef CONFIG_H
#define CONFIG_H

const long serialBaudRate = 115200;
const int micPin = 35;
const int numberMeasures = 100;

// WiFi config
char wifi_ssid[] = "MIWIFI_mnQN";
char wifi_password[] = "jrcCHfac";
const char test_server[] = "192.168.1.227";
const int test_connection_port = 22;

// SNTP config
const int NTP_PACKET_SIZE = 48;
const int udpPort = 2390;
const int ntpPort = 123;
const int daylightOffset_sec = 3600;

// Function prototypes
void beginUdp();
void connectWifi();
void readMeasure();
long fetchTimeNtp();

#endif
