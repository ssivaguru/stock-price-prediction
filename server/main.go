package main

import (
	"log"
	apiserver "stock-price-prediction/predictionServer/apiServer"
)

func main() {
	exit := make(chan bool)

	server := &ServerStruct{ApiServer: apiserver.New()}

	go server.StartServer()

	select {
	case <-exit:
		log.Println("closing application")
	}
}
