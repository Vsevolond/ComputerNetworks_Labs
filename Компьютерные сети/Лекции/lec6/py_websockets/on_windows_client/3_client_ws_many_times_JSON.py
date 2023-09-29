#!/usr/bin/python3

import asyncio
import websockets

async def hello():
    #uri = "wss://echo.websocket.events"
    uri = "ws://151.248.113.144:8000"
    async with websockets.connect(uri) as websocket:

        for i in range(0,10):
            await websocket.send('{"a":{"x":"1","y":"2"},"b":{"x":"1","y":"2"}}')
            name = await websocket.recv()
            print("Reply from WS"+str(i)+":"+str(name))


asyncio.get_event_loop().run_until_complete(hello())
