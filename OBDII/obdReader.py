import obd
import paho.mqtt.client as mqtt
import time
client =mqtt.Client()
client.username_pw_set("csc591","iot123!@#")
client.connect("35.196.98.6",1883,60)
'''
while True:
    msg="60,30,800"
    client.publish("topic/obd2",msg)
    print(msg+" was published")
    time.sleep(5)
'''
obd.logger.setLevel(obd.logging.DEBUG)

ports =obd.scan_serial()
print("ports:")
print ports

connection=obd.OBD(ports[0])
print("connection status:")
print(connection.status())

commands = connection.supported_commands
for command in commands:
    print(command.name)
while True:
    #command = input("enter command(type 'quit' to exit):")
    #if (command=="quit"):
     #   break;
    try:
        speed=connection.query(obd.commands.SPEED)   #kmph
        print(speed.value)
       # fuel=connection.query(obd.commands.FUEL_LEVEL) # unit.percent
       # print(fuel.value)
        rpm=connection.query(obd.commands.RPM)       #unit.rpm
        print(rpm.value)
        speed = str(speed).split()[0]
        rpm= str(rpm).split()[0]
       # fuel=str(fuel).split()[0]
        fuel='45'
        if speed != 'None' and rpm != 'None' and fuel != 'None':
            msg =speed + ","+fuel +","+rpm
            client.publish("topic/obd2",msg)
            print(msg+" was published")
        time.sleep(5)
    except Exception as ex:
        print("error" + str(ex))
connection.close()
client.disconnect()
