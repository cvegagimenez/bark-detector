/*
File where we will configure the SNTP service to get the exact time.
*/

#include <WiFiUdp.h>
#include <precise_sntp.h>
#include "config.h"


unsigned int localPort = udpPort;
char timeServer[] = "pool.ntp.org";
byte packetBuffer[ NTP_PACKET_SIZE ];

WiFiUDP udp;
precise_sntp sntp(udp, timeServer);

/* 
  Function to start SNTP.
  @return void
*/
void beginUdp(){
    udp.begin(localPort);
    Serial.println("Started UDP");
}

/*
  Function to parse and calculate the time that we got from the ntp Server.
  @return unsigned long
*/

long fetchTimeNtp(){
  sntp.update();
  return sntp.get_epoch();
}

