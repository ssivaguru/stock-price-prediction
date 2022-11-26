#!/usr/bin/env python

import asyncio
import websockets

def handleMessage(message):
    splitString = message.split()
    resMessage = ""
    if splitString[0] == "predict":
        ##call predictor class
        print("check with predictor")
        resMessage = "Handled Message"
    else:    
        print("Unhandled message")
        resMessage = "Unhandled message"
    
    return resMessage

async def handler(websocket):
    while True:
        try:
            message = await websocket.recv()
        except websockets.ConnectionClosedOK:
            print("websocket connection closed")
            break

        resp = handleMessage(message)

        await websocket.send(resp)


async def main():
    async with websockets.serve(handler, "", 8000):
        await asyncio.Future()  # run forever


if __name__ == "__main__":
    try:
        asyncio.run(main())
        print("test")
    except KeyboardInterrupt:
        print("error")