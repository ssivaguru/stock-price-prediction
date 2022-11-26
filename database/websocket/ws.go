package websocket

import (
	"log"
	"net/http"

	websocket "github.com/gorilla/websocket"
)

type WsInterface interface {
	StartServer(connCh chan ConnHandler) error
	Close()
}

type WsStruct struct {
	isClose chan bool
}

func New() WsInterface {
	return &WsStruct{isClose: make(chan bool)}
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ws *WsStruct) wsEndpoint(connCh chan ConnHandler, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	//this will almost be instatinous
	connCh <- CreateNewConn(conn)
}

func (ws *WsStruct) setupRoutes(connCh chan ConnHandler) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.wsEndpoint(connCh, w, r)
	})
}

func (ws *WsStruct) StartServer(connCh chan ConnHandler) error {
	ws.setupRoutes(connCh)
	log.Println("starting local host")
	return http.ListenAndServe(":8080", nil)
}

func (ws *WsStruct) Close() {

}
