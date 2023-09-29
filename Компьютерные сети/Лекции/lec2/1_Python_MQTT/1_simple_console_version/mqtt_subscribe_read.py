import paho.mqtt.client as mqtt

def on_connect(client, userdata, flags, rc):
    print("Connected with result code "+str(rc))
    client.subscribe("IU/9")

def on_message(client, userdata, msg):
    if msg.topic == 'IU/9':
        print(msg.topic+" "+str(msg.payload.decode('UTF-8')))

mqtt_host = 'test.mosquitto.org'
mqtt_port = 1883

topic = 'IU/9'
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message
client.connect(mqtt_host, mqtt_port)
client.loop_start()



