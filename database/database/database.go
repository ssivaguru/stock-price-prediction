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
	UpdateByName(name string, predict []float64) (*mongo.UpdateResult, error)
	UpdateState(name string, state int) (*mongo.UpdateResult, error)
}

type DbStruct struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type Stock struct {
	ID              primitive.ObjectID `bson:"_id"`
	Name            string             `bson:"name"`
	Path            string             `bson:"path"`
	Status          uint               `bson:"status"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	Requesting_user []string           `bson:"requesting_user"`
	Predicted_val   []float64          `bson:"predicted_val"`
}

func ConnectDb() DbInterface {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println(err)
		return nil
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println(err)
		return nil
	}

	collection := client.Database("stockPredictor").Collection("stocks")

	server := &DbStruct{client, collection}

	return server
}

func (db *DbStruct) InsertData(task *Stock) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := db.collection.InsertOne(ctx, task)
	return err
}

func (db *DbStruct) GetData(name string) (*Stock, error) {
	filter := bson.D{
		primitive.E{Key: "name", Value: name},
	}

	return db.fileterData(filter)
}

func (db *DbStruct) UpdateState(name string, state int) (*mongo.UpdateResult, error) {
	filter := bson.D{
		primitive.E{Key: "name", Value: name},
	}

	update := bson.D{
		primitive.E{"$set",
			bson.D{
				primitive.E{Key: "status", Value: state},
			},
		}}

	result, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	return result, nil
}

func (db *DbStruct) UpdateByName(name string, prediction []float64) (*mongo.UpdateResult, error) {

	filter := bson.D{
		primitive.E{Key: "name", Value: name},
	}

	update := bson.D{
		primitive.E{"$set",
			bson.D{
				primitive.E{Key: "predicted_val", Value: prediction},
			},
		}}

	result, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	return result, nil
}

func (db *DbStruct) DeleteStock(text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
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
