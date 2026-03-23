#include "Arduino.h"
#include "config.h"

void setup() {
  Serial.begin(serialBaudRate);
  connectWifi();
  beginUdp();
  connectMqtt();
}

void loop() {
  ensureMqttConnected();
  readMeasure();
  delay(5);
}
