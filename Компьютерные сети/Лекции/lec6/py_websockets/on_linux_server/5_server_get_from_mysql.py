#!/usr/bin/python3
import datetime
import time
import asyncio
import websockets
import pymysql


async def hello(websocket, path):
    name = await websocket.recv()
    print(str(name))

    greeting = "Hello "+str(name)+ " at " + str(datetime.datetime.now())+"!<br>"

    con = pymysql.connect('185.4.72.67','testizo','123','testizo')
    cur = con.cursor()
    cur.execute("SELECT * FROM _epi_danila")
    rows = cur.fetchall()

    for row in rows:
        greeting = greeting + str(row[0])+" "+str(row[1])+" "+str(row[2])+"<br>"
    con.close()


    await websocket.send(greeting)
    print(str(greeting))





start_server = websockets.serve(hello, "151.248.113.144", 8000)

asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
