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

