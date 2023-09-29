#!/usr/bin/python3
import socket

a = input()

sock = socket.socket()
sock.connect(('151.248.113.144', 9090))

for i in range(1,1000):
    b = str(i)+a
    sock.send(b.encode("utf-8"))
    data = sock.recv(1024)
    print(data.decode("utf-8"))
sock.close()

