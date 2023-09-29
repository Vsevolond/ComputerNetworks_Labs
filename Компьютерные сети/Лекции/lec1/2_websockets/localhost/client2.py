#!/usr/bin/python3


import asyncio
import websockets

async def hello():
    uri = "ws://localhost:8765"
    async with websockets.connect(uri) as websocket:
        name = input("What's your name? ")

        await websocket.send(name)
        print(str(name))

        greeting = await websocket.recv()
        print(str(greeting))

asyncio.get_event_loop().run_until_complete(hello())
