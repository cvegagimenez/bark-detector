/*
  In this file we will configure all the communications needed to establish the connectivity between
  the embedded system and the MQTT broker.
*/

#include <WiFi.h>
#include <WiFiClient.h>
#include <PubSubClient.h>
#include "config.h"


WiFiClient wifi_client;
PubSubClient mqtt_client(wifi_client);

void connectWifi(){
  WiFi.begin(wifi_ssid, wifi_password);
  
  while ( WiFi.status() != WL_CONNECTED) {
    delay(300);
  }
  
  while (WiFi.localIP() == INADDR_NONE) {
    delay(300);
  }
  
  Serial.print("WiFi connected: ");
  Serial.println(WiFi.localIP());
}

void connectMqtt() {
  mqtt_client.setServer(mqtt_broker, mqtt_port);

  uint8_t tries = 0;
  while (!mqtt_client.connected()) {
    if (mqtt_client.connect(mqtt_client_id)) {
      Serial.println("MQTT connected");
      return;
    }

    if (tries++ > 100) {
      Serial.println("MQTT broker not responding");
      while(1);
    }
    delay(100);
  }
}

void ensureMqttConnected() {
  if (WiFi.status() != WL_CONNECTED) {
    connectWifi();
  }

  if (!mqtt_client.connected()) {
    connectMqtt();
  }

  mqtt_client.loop();
}

bool publishMetric(long epochTime, float rms) {
  ensureMqttConnected();

  String payload = String(epochTime) + "|" + String(sensor_id) + "|" + String(rms, 4);
  bool published = mqtt_client.publish(mqtt_topic, payload.c_str());

  if (!published) {
    Serial.println("Failed to publish metric to MQTT");
    return false;
  }

  return true;
}
