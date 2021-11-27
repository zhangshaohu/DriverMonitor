import paho.mqtt.client as mqtt
import datetime as dt
import os
from influxdb import InfluxDBClient

def on_connect(client, userdata, flags, rc):
    print(f"Connected with result code {rc}")

    client.subscribe("topic/watch")
    client.subscribe("topic/obd2")


def mqtt_to_influx_point(payload, topic):
    base = [{
        "measurement": "records",
        "tags": {
            "Token": "test",
            },
        "time": dt.datetime.now(),
        "fields": {}
    }]

    data = payload.decode('utf-8').split(',')
    if topic == "topic/watch":
        # Smartwatch payload
        base[0]["fields"]["Lat"] = float(data[0])
        base[0]["fields"]["Lng"] = float(data[1])
        base[0]["fields"]["HeartRate"] = float(data[2])
    elif topic == "topic/obd2":
        # OBD-II sensor payload
        base[0]["fields"]["Speed"] = float(data[0])
        base[0]["fields"]["FuelRemaining"] = float(data[1])
        base[0]["fields"]["RPM"] = float(data[2])

    return base

def on_message(client, userdata, msg):
    influx_client.write_points(mqtt_to_influx_point(msg.payload, msg.topic))

if __name__ == '__main__':
    influx_client = InfluxDBClient("influxdb", 8086, database="csc591")

    client = mqtt.Client()
    client.username_pw_set(os.environ["MQTT_USER"], password=os.environ["MQTT_PASSWORD"])
    client.on_connect = on_connect
    client.on_message = on_message

    client.connect(os.environ["MQTT_HOST"])

    client.loop_forever()
