package producer

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "train-stock"
	brokerAddress = "localhost:9092"
)

type ProducerInterface interface {
	PublishMessage(string, []byte) bool
}

type producerStruct struct {
	w *kafka.Writer
}

func New() ProducerInterface {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		// assign the logger to the writer
		Logger: log.New(os.Stdout, "Producer:: ", 1),
	})

	return &producerStruct{w}
}

func (pub *producerStruct) PublishMessage(topic string, msg []byte) bool {

	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	err := pub.w.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   []byte("Name"),
		// create an arbitrary message payload for the value
		Value: msg,
	})
	if err != nil {
		log.Println("produce:: could not write message " + err.Error())
		return false
	}

	return true
}
