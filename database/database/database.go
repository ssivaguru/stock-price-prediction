package database

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DbInterface interface {
	InsertData(task *Stock) error
	GetData(name string) (*Stock, error)
	DeleteStock(text string) error
}

type DbStruct struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type Stock struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Path      string             `bson:"path"`
	Status    uint               `bson:"status"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

var ctx = context.TODO()

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

	collection := client.Database("stockPredictor").Collection("stocks")

	server := &DbStruct{client, collection}

	return server
}

func (db *DbStruct) InsertData(task *Stock) error {
	_, err := db.collection.InsertOne(ctx, task)
	return err
}

func (db *DbStruct) GetData(name string) (*Stock, error) {
	filter := bson.D{
		primitive.E{Key: "name", Value: name},
	}

	return db.fileterData(filter)
}

func (db *DbStruct) DeleteStock(text string) error {
	filter := bson.D{primitive.E{Key: "name", Value: text}}

	res, err := db.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("No tasks were deleted")
	}

	return nil
}

func (db *DbStruct) fileterData(filter interface{}) (*Stock, error) {
	// A slice of tasks for storing the decoded documents
	var tasks *Stock

	cur, err := db.collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}

	// Iterate through the cursor and decode each document one at a time
	for cur.Next(ctx) {
		var t Stock
		err := cur.Decode(&t)
		if err != nil {
			return tasks, err
		}

		tasks = &t
		break
	}

	if err := cur.Err(); err != nil {
		return tasks, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if tasks == nil {
		return tasks, mongo.ErrNoDocuments
	}

	return tasks, nil
}
