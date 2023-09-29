#!/usr/bin/python3

import asyncio
import websockets

async def hello():
    uri = "ws://151.248.113.144:8765"
    async with websockets.connect(uri) as websocket:
        await websocket.send("Hello world!")
        await websocket.recv()

asyncio.get_event_loop().run_until_complete(hello())
