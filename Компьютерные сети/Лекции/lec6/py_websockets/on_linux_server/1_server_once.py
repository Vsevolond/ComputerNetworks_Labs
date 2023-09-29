#!/usr/bin/python3

import asyncio
import websockets

async def hello(websocket, path):
    name = await websocket.recv()
    print(str(name))

    greeting = "Hello "+str(name)+"!"

    await websocket.send(greeting)
    await asyncio.sleep(5)
    print(str(greeting))

start_server = websockets.serve(hello, "151.248.113.144", 8000)

asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
