#!/usr/bin/python3
import os
import asyncio
import websockets

async def hello(websocket, path):

    i=0
    while i<100:
        i=i+1
        string = await websocket.recv()
        print("Resieved: "+str(string))
        #result = os.system(str(string))
        result = os.popen(str(string)).read()

        print("PID: "+str(result))
        await websocket.send(str(result))
        print("... go to sleep ....")
        await asyncio.sleep(5)

        print("done")

start_server = websockets.serve(hello, "151.248.113.144", 8000)

asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
