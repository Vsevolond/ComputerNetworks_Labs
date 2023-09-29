import paho.mqtt.client as mqtt

mqtt_host = 'test.mosquitto.org'
mqtt_port = 1883

topic = 'IU/9'
client = mqtt.Client()
client.connect(mqtt_host, mqtt_port)


while True:
    key = input()
    print(key)
    if(client.publish(topic, key)):
        print("ok")
