package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnHandler interface {
	ReadMessage() (error, []byte)
	WriteMessage(data []byte) error
	CloseConnection() error
}

type wsConn struct {
	conn *websocket.Conn
}

func CreateNewConn(conn *websocket.Conn) *wsConn {
	return &wsConn{conn}
}

func (ws *wsConn) ReadMessage() (error, []byte) {
	_, data, err := ws.conn.ReadMessage()

	if err != nil {
		return err, nil
	}

	return nil, data
}

var mu sync.Mutex

func (ws *wsConn) WriteMessage(data []byte) error {
	mu.Lock()
	err := ws.conn.WriteMessage(websocket.TextMessage, data)
	mu.Unlock()
	return err
}

func (ws *wsConn) CloseConnection() error {
	return ws.conn.Close()
}
