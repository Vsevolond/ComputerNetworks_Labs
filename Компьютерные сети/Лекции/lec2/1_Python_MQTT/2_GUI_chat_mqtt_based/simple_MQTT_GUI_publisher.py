import paho.mqtt.client as mqtt
import PySimpleGUI as sg

layout = [
    [sg.Text('Message'), sg.InputText()],
    [sg.Output(size=(88, 20))],
    [sg.Submit(), sg.Cancel()]
]

window = sg.Window('MQTT chat', layout)

mqtt_host = 'test.mosquitto.org'
mqtt_port = 1883

topic = 'IU/9'
client = mqtt.Client()
client.connect(mqtt_host, mqtt_port)

while True:                             
    event, values = window.read()
    # print(event, values) #debug
    if event in (None, 'Cancel'):
        break
    if event == 'Submit':
        # print("amount of fields = "+str(len(values)))
        if values[0]:
            if(client.publish(topic, str(values[0]))):
                print("sent to topic: "+str(values[0]))
