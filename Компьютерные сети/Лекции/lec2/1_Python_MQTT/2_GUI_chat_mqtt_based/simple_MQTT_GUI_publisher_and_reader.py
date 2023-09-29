import paho.mqtt.client as mqtt
import PySimpleGUI as sg

def on_connect(client, userdata, flags, rc):
    print("Connected with result code "+str(rc))
    client.subscribe("IU/9")

def on_message(client, userdata, msg):
    if msg.topic == 'IU/9':
        print(msg.topic+" "+str(msg.payload.decode('UTF-8')))

layout = [
    [sg.Text('Name'), sg.InputText()],
    [sg.Text('Message'), sg.InputText()],
    [sg.Output(size=(88, 20))],
    [sg.Submit(), sg.Cancel()]
]

window = sg.Window('MQTT chat', layout)

mqtt_host = 'test.mosquitto.org'
mqtt_port = 1883

topic = 'IU/9'
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message
client.connect(mqtt_host, mqtt_port)

while True:                             
    event, values = window.read()
    # print(event, values) #debug
    if event in (None, 'Cancel'):
        break
    if event == 'Submit':
        # print("amount of fields = "+str(len(values)))
        if values[0]:
            if(client.publish(topic, str(values[0])+ "|---say ---> " + str(values[1]))):
                print("sent to topic: "+str(values[1]))
                client.loop_start()
