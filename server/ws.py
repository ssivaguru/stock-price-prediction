import asyncio
import websockets

class WebsocketServer:
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

    async def handler(self, websocket):
        while True:
            try:
                message = await websocket.recv()
            except websockets.ConnectionClosedOK:
                print("websocket connection closed")
                break

            resp = self.handleMessage(message)

            await websocket.send(resp)