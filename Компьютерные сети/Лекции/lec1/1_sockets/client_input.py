#!/usr/bin/python3
import socket

a = input()

sock = socket.socket()
sock.connect(('151.248.113.144', 9090))
sock.send(a.encode("utf-8"))

data = sock.recv(1024)
sock.close()
print(data.decode("utf-8"))
