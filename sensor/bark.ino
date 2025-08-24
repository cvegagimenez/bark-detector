#include <Arduino.h>

const int micPin = 35;
const int numberMeasures = 100;
int samples[numberMeasures];
int idx = 0;

void setup() {
  Serial.begin(115200);
}

void loop() {
  int raw = analogRead(micPin);

  static float dc = 1500;
  dc += 0.001 * (raw - dc);

  samples[idx] = raw - dc;

  idx++;
  if (idx >= numberMeasures) {
    idx = 0;
    Serial.println(calculateRMS(samples));
  }

  delay(5);
}

float calculateRMS(int measures[]) {
  long sumSq = 0;
  for (int i = 0; i < numberMeasures; i++) {
    sumSq += (long)measures[i] * (long)measures[i];
  }
  return sqrt((float)sumSq / numberMeasures);
}