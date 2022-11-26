package main

import (
	"log"
	"stock-price-prediction/database/database"
	"stock-price-prediction/database/websocket"
)

func main() {

	exit := make(chan bool)
	server := &DbServer{Dbclient: database.ConnectDb(), connError: make(chan error), WsServer: websocket.New()}
	go server.start(exit)

	select {
	case <-exit:
		log.Println("closing application")
	}
}
