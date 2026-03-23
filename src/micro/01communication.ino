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
    Serial.println("Connecting to WiFi");
    delay(300);
  }
  
  Serial.println("\nYou're connected to the network");
  Serial.println("Waiting for an ip address");
  
  while (WiFi.localIP() == INADDR_NONE) {
    Serial.print(".");
    delay(300);
  }
  
  Serial.println("\nIP Address obtained");
  Serial.println(WiFi.localIP());
}

void connectMqtt() {
  mqtt_client.setServer(mqtt_broker, mqtt_port);

  uint8_t tries = 0;
  while (!mqtt_client.connected()) {
    Serial.println("Connecting to MQTT broker");
    if (mqtt_client.connect(mqtt_client_id)) {
      Serial.println("Connected to the MQTT broker");
      return;
    }

    Serial.print(".");
    if (tries++ > 100) {
      Serial.println("\nThe MQTT broker isn't responding");
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

  Serial.println("Published metric to MQTT: " + payload);
  return true;
}
