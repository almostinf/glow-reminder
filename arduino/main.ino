#include <ESP8266WiFi.h>
#include <ESP8266WebServer.h>
#include <ArduinoJson.h>


const char* ssid = "YOUR WIFI";
const char* password = "WIFI PASSWORD";

ESP8266WebServer server(80);

uint8_t LED1pin = D8;  // Пин для цвета 3
uint8_t LED2pin = D6;  // Пин для цвета 2
uint8_t LED3pin = D5; // Пин для цвета 1

void setup() 
{
  Serial.begin(115200);
  pinMode(LED1pin, OUTPUT);
  pinMode(LED2pin, OUTPUT);
  pinMode(LED3pin, OUTPUT);

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.println("Connecting to WiFi...");
  }

  Serial.println("Connected to WiFi");
  Serial.println(WiFi.localIP());
  
  server.on("/glow_reminder", HTTP_POST, handle_glow_reminder);
  server.onNotFound(handle_NotFound);
  
  server.begin();
  Serial.println("HTTP server started");
}

void loop() 
{
  server.handleClient();
}

void handle_glow_reminder() 
{
  if (server.hasArg("plain"))
  {
    String body = server.arg("plain");
    DynamicJsonDocument doc(1024);
    deserializeJson(doc, body);

    int colour = doc["colour"];
    int mode = doc["mode"];

    server.send(200, "text/plain", "OK");

    if (colour == 1) {
      controlLED(LED3pin, mode);
    }
    else if (colour == 2) 
    {
      controlLED(LED2pin, mode);
    } 
    else if (colour == 3) 
    {
      controlLED(LED1pin, mode);
    }
  } 
  else 
  {
    server.send(400, "text/plain", "Bad Request");
  }
}

void controlLED(uint8_t pin, int mode) 
{
  if (mode == 1) 
  {
    digitalWrite(pin, HIGH);
    delay(10000);  // 10 секунд
    digitalWrite(pin, LOW);
  } 
  else if (mode == 2) 
  {
    for (int i = 0; i < 10; i++) 
    {
      digitalWrite(pin, HIGH);
      delay(500);  // 0.5 секунд
      digitalWrite(pin, LOW);
      delay(500);  // 0.5 секунд
    }
  }
}

void handle_NotFound()
{
  server.send(404, "text/plain", "Not found");
}