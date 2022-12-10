package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "localhost:9092"
)

func main() {
	//trainData("MARUTI.NS", nil)
	//return
	trainCh := make(chan int, 10)
	predictCh := make(chan int, 50)
	for true {
		fmt.Println("Consuming topic")
		consumePrediction(predictCh)
		time.Sleep(time.Second * 10)
		consumeTrain(trainCh)
	}
}

func consumePrediction(predict chan int) {

	fmt.Println("started dailer")

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       &tls.Config{},
	}

	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   "predict",
		GroupID: "my-group",
		Dialer:  dialer,
	})

	fmt.Println("reader complete")
	count := 0
	//break after 30 tries
	for {
		if count == 30 {
			return
		}
		//channel is full
		if len(predict) == cap(predict) {
			return
		}
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("could not read message " + err.Error())
			count++
			continue
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}

func consumeTrain(train chan int) {

	fmt.Println("started dailer")

	dialer := &kafka.Dialer{
		Timeout:   5 * time.Second,
		DualStack: true,
		TLS:       &tls.Config{},
	}

	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   "train",
		GroupID: "my-group",
		Dialer:  dialer,
	})
	fmt.Println("reader complete")
	count := 0
	//break after 10 tries
	for {
		if count == 10 {
			return
		}

		//channel is full
		if len(train) == cap(train) {
			return
		}
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("could not read message " + err.Error())
			count++
			continue
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}

func GetStringInBetween(str string, start string, end string) (result float64) {
	fmt.Println(str)
	res, _ := strconv.ParseFloat(str[1:len(str)-1], 64)
	return res
}

func trainData(msg string, close chan int) {
	//Assuming that the data has been trained
	var b2 bytes.Buffer
	var b3 bytes.Buffer

	cmd := exec.Command("python", "/home/siva/Github/src/stock-price-prediction/trainer/ml.py", "predict", msg)
	cmd.Stdout = &b2
	cmd.Stderr = &b3
	err := cmd.Run()

	if err != nil {
		println(err.Error())
		return
	}

	splitString := strings.Split(b2.String(), "\n")

	found := false
	var predict []float64
	for _, val := range splitString {
		if strings.Contains(val, ":OUT:") {
			found = true
			continue
		}
		if found && len(strings.TrimSpace(val)) != 0 {
			predict = append(predict, GetStringInBetween(val, "[", "}"))
		}
	}

	resp := make(map[string]interface{})
	resp["name"] = msg
	resp["predict"] = predict
	out, err := json.Marshal(resp)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(out[:]))

	client := &http.Client{}
	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/stock", bytes.NewBuffer(out))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.StatusCode)

	//<-close
}

func predictData(msg string, close chan int) {

	<-close
}
