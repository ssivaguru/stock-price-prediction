import websocket
import _thread


class WS:
    isRunning = True

    def ws_message(ws, message):
        print("WebSocket thread: %s" % message)

    def ws_open(ws):
        print("websocket connected")

    def ws_error(ws, err):
        print("error recevied ", err)

    def we_close(ws, status, msg):
        print("websocket closed ", status, msg)
    
    def send_msg(self, msg):
        self.ws.send(msg)
    
    def ws_thread(self, *args):
        self.isRunning = True
        websocket.setdefaulttimeout(3)
        self.ws = websocket.WebSocketApp("ws://localhost:8000", on_open = self.ws_open, on_message = self.ws_message)
        self.ws.run_forever()
        self.isRunning = False

    def RunAsync(self):
        # Start a new thread for the WebSocket interface
        _thread.start_new_thread(self.ws_thread, ())















