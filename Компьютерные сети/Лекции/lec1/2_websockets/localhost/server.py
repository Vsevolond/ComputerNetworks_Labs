#!/usr/bin/python3

import asyncio
import websockets

async def hello(websocket, path):
    name = await websocket.recv()
    print(str(name))

    greeting = "Hello "+str(name)+"!"

    await websocket.send(greeting)
    print(str(greeting))

start_server = websockets.serve(hello, "localhost", 8765)

asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
