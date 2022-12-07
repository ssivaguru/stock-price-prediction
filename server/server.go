package main

import (
	apiserver "stock-price-prediction/predictionServer/apiServer"
)

type ServerStruct struct {
	ApiServer apiserver.ApiInterface
	connError chan error
}

func (s *ServerStruct) StartServer() {
	err := s.ApiServer.StartServer()
	s.connError <- err
}
