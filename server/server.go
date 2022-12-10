package main

import (
	apiserver "stock-price-prediction/predictionServer/apiServer"
)

type ServerStruct struct {
	ApiServer apiserver.ApiInterface
	connError chan error
}

func (s *ServerStruct) StartServer() {
	msgCh := make(chan string, 1)
	err := s.ApiServer.StartServer(msgCh)
	s.connError <- err
}
