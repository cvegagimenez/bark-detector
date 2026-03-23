/*
  In this file we will configure all the communications needed to establish the connectivity between
  the embedded system and the MQTT broker.
*/

#include <WiFi.h>
#include <WiFiClient.h>
#include "config.h"


WiFiClient wifi_client;

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
  
  uint8_t tries = 0;
  while (wifi_client.connect(test_server, test_connection_port) == false) {
    Serial.print(".");
    if (tries++ > 100) {
      Serial.println("\nThe server isn't responding");
      while(1);
    }
    delay(100);
  }
  Serial.println("\nConnected to the server!");
}



