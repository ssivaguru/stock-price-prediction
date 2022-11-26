package main

import (
	"stock-price-prediction/database/database"
)

func main() {

	server := &DbServer{Dbclient: database.ConnectDb()}
	server.start()

}
