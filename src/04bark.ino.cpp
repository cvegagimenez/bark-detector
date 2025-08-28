# 1 "C:\\Users\\Gigi\\AppData\\Local\\Temp\\tmpkmovjrz2"
#include <Arduino.h>
# 1 "D:/Workspace/bark-detector/src/04bark.ino"
#include "Arduino.h"
#include "config.h"
void setup();
void loop();
void connectWifi();
void beginUdp();
long fetchTimeNtp();
void readMeasure();
float calculateRMS(int measures[]);
#line 4 "D:/Workspace/bark-detector/src/04bark.ino"
void setup() {
  Serial.begin(serialBaudRate);
  connectWifi();
  beginUdp();
}

void loop() {
  readMeasure();
  delay(5);
}
# 1 "D:/Workspace/bark-detector/src/01communication.ino"





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
# 1 "D:/Workspace/bark-detector/src/02sntp.ino"




#include <WiFiUdp.h>
#include <precise_sntp.h>
#include "config.h"


unsigned int localPort = udpPort;
char timeServer[] = "pool.ntp.org";
byte packetBuffer[ NTP_PACKET_SIZE ];

WiFiUDP udp;
precise_sntp sntp(udp, timeServer);





void beginUdp(){
    udp.begin(localPort);
    Serial.println("Started UDP");
}






long fetchTimeNtp(){
  sntp.update();
  return sntp.get_epoch();
}
# 1 "D:/Workspace/bark-detector/src/03audio.ino"
#include "Arduino.h"
#include "config.h"

int samples[numberMeasures];
int idx = 0;


void readMeasure(){
  int raw = analogRead(micPin);

  static float dc = 1500;
  dc += 0.001 * (raw - dc);

  samples[idx] = raw - dc;

  idx++;
  if (idx >= numberMeasures) {
    idx = 0;
    float rms = calculateRMS(samples);
    long epochTime = fetchTimeNtp();
    Serial.println(String(epochTime) + " | " + String(rms));
  }
}

float calculateRMS(int measures[]) {
  long sumSq = 0;
  for (int i = 0; i < numberMeasures; i++) {
    sumSq += (long)measures[i] * (long)measures[i];
  }
  return sqrt((float)sumSq / numberMeasures);
}