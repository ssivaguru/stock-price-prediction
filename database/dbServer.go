package main

import (
	"log"
	"stock-price-prediction/database/database"
	"stock-price-prediction/database/websocket"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
				log.Println("New Connection")
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

WS_LOOP:
	for {
		err, msg := conn.ReadMessage()

		if err != nil {
			log.Println("handleConnection:: error while reading message ", err)
			break WS_LOOP
		}

		if msg == nil {
			continue
		}

		go s.handleMessage(msg, conn)
	}

}

func (s *DbServer) handlePredict(msg string) string {
	name := strings.Split(msg, " ")[1]

	_, err := s.Dbclient.GetData(name)

	if err == mongo.ErrNoDocuments {
		return "Not Created"
	} else if err != nil {
		log.Println("handlePredict:: error getting stock name ", err)
		return "error getting stock name"
	}

	//we need to do more handling
	return "Yes"
}

func (s *DbServer) handleCreateStock(msg string) string {
	name := strings.Split(msg, " ")[1]

	//jsut check if its already created

	_, err := s.Dbclient.GetData(name)

	if err == nil {
		//then stock is alreadty created
		return "Stock is already created"
	}

	err = s.Dbclient.InsertData(&database.Stock{ID: primitive.NewObjectIDFromTimestamp(time.Now()), Name: name, Path: "", Status: StatusTraning, UpdatedAt: time.Now()})

	if err != nil {
		log.Println("handleCreateStock:: error creating stock ", err)
		return "error creating stock"
	}

	return "Stock Created"
}

func (s *DbServer) handleDescribe(msg string) string {
	return ""
}

func (s *DbServer) handleDeleteStock(msg string) string {
	return ""
}

func (s *DbServer) handleUpdateStock(msg string) string {
	return ""
}

func (s *DbServer) handleMessage(msg []byte, conn websocket.ConnHandler) {

	strMsg := string(msg)
	resp := "No Resp"
	if strings.HasPrefix(strMsg, PredictStock) {
		resp = s.handlePredict(strMsg)
	} else if strings.HasPrefix(strMsg, CreateStock) {
		resp = s.handleCreateStock(strMsg)
	} else if strings.HasPrefix(strMsg, Describe) {
		resp = s.handleDescribe(strMsg)
	} else if strings.HasPrefix(strMsg, DeleteStock) {
		resp = s.handleDeleteStock(strMsg)
	} else if strings.HasPrefix(strMsg, UpdateStock) {
		resp = s.handleUpdateStock(strMsg)
	} else {
		log.Println("handleMessage:: unhandled message")
	}

	conn.WriteMessage([]byte(resp))
}

func (s *DbServer) writeMessage() {

}
