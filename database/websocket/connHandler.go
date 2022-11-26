package websocket

import "github.com/gorilla/websocket"

type ConnHandler interface {
	ReadMessage() (error, string)
	WriteMessage(data []byte) error
	CloseConnection() error
}

type wsConn struct {
	conn *websocket.Conn
}

func CreateNewConn(conn *websocket.Conn) *wsConn {
	return &wsConn{conn}
}

func (ws *wsConn) ReadMessage() (error, string) {
	_, data, err := ws.conn.ReadMessage()

	if err != nil {
		return err, ""
	}

	return nil, string(data[:])
}

func (ws *wsConn) WriteMessage(data []byte) error {
	return ws.conn.WriteMessage(websocket.TextMessage, data)
}

func (ws *wsConn) CloseConnection() error {
	return ws.conn.Close()
}
