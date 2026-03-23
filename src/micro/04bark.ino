#include "Arduino.h"
#include "config.h"

void setup() {
  Serial.begin(serialBaudRate);
  connectWifi();
  beginUdp();
}

void loop() {
  readMeasure();
  delay(5);
}
