package main

import (
	"log"
	"stock-price-prediction/database/database"
	"stock-price-prediction/database/websocket"
	"time"
)

type DbServer struct {
	Dbclient  database.DbInterface
	WsServer  websocket.WsInterface
	connError chan error
}

func (s *DbServer) startServer(connCh chan websocket.ConnHandler) {
	err := s.WsServer.StartServer(connCh)
	s.connError <- err
}

func (s *DbServer) start(exit chan bool) {

	wsConn := make(chan websocket.ConnHandler)
	go s.startServer(wsConn)

MAIN_LOOP:
	for {
		select {
		case conn := <-wsConn:

			if conn != nil {
				go s.handleConnection(conn)
			}
		case err := <-s.connError:
			log.Println("main::start:: server connection error ", err)
			break MAIN_LOOP
		}
	}

	exit <- true
}

func (s *DbServer) handleConnection(conn websocket.ConnHandler) {
	err, message := conn.ReadMessage()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(message)

	conn.WriteMessage([]byte("Yolo Polo"))
	time.Sleep(time.Second * 10)
	conn.CloseConnection()
}
