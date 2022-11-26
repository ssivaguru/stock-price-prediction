package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DbInterface interface {
}

type DbStruct struct {
	client *mongo.Client
}

func ConnectDb() DbInterface {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println(err)
		return nil
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println(err)
		return nil
	}
	server := &DbStruct{client}

	return server
}
