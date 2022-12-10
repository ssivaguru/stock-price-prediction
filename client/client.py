from ws import WS
import sys
import time
import websocket

if len(sys.argv) != 2:
    print("stock name needs to be passed as command line argument")
    #exit(1)

websocket = WS()

websocket.RunAsync()

while websocket.isRunning:
    time.sleep(5)
    print("Main thread: %d" % time.time())