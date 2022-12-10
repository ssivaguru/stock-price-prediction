package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	producer "stock-price-prediction/predictionServer/apiServer/kafka-producer"

	"github.com/gin-gonic/gin"
)

type ApiInterface interface {
	StartServer(msgCh chan string) error
	Close()
}

type ApiStruct struct {
	isClose   chan bool
	PubServer producer.ProducerInterface
}

func New() ApiInterface {
	return &ApiStruct{isClose: make(chan bool), PubServer: producer.New()}
}

func getJson(data map[string]string) []byte {
	jsonData, _ := json.Marshal(data)
	log.Println(string(jsonData[:]))
	return jsonData
}

func (ws *ApiStruct) handlePredit(c *gin.Context) {
	paramPairs := c.Request.URL.Query()
	if _, ok := paramPairs["name"]; !ok {
		c.JSON(400, "{Msg: Malformed URL}")
		return
	}

	resp, err := http.Get("http://localhost:8080/api/predict?name=" + paramPairs.Get("name"))
	if err != nil {
		log.Println(err)
		c.Data(500, gin.MIMEJSON, getJson(map[string]string{"error": "DB server down"}))
		return
	}

	var data map[string]interface{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		c.Data(500, gin.MIMEJSON, getJson(map[string]string{"error": "error parsing json body"}))
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		c.Data(500, gin.MIMEJSON, getJson(map[string]string{"error": "error getting json data"}))
		return
	}
	//then this is new data start prediction
	if _, ok := data["Status"]; ok {
		if data["Status"] == 1 {
			c.Data(201, gin.MIMEJSON, getJson(map[string]string{"status": "traning the data model"}))
			ws.PubServer.PublishMessage("train", body)
			return
			//then we have a trained data
		} else if data["Status"] == 3 {
			c.Data(200, gin.MIMEJSON, getJson(map[string]string{"status": "model is ready so lets start predicting"}))
			ws.PubServer.PublishMessage("predict", body)
			return
		}

	}

	//means we have to push to kafka
	if resp.StatusCode == 201 {
		c.Data(200, gin.MIMEJSON, getJson(map[string]string{"status": "model is ready so lets start predicting"}))
		if ws.PubServer.PublishMessage("train", body) {
			log.Println("data has been ploted")
		} else {
			log.Println("data has not been ploted")
		}
		return
	}
	c.JSON(resp.StatusCode, data)
}

func handleDescribe(c *gin.Context) {
	paramPairs := c.Request.URL.Query()
	for key, values := range paramPairs {
		fmt.Printf("key = %v, value(s) = %v\n", key, values)
	}
	c.JSON(200, "{Status:Success}")
}

func (ws *ApiStruct) setupRoutes(router *gin.Engine, msgCh chan string) {
	router.GET("/api/predict", func(c *gin.Context) {
		ws.handlePredit(c)
	})

	router.GET("/api/describe", func(c *gin.Context) {
		handleDescribe(c)
	})
}

func (ws *ApiStruct) StartServer(msgCh chan string) error {
	ws.PubServer.PublishMessage("train", []byte("Test"))
	return nil
	router := gin.Default()
	ws.setupRoutes(router, msgCh)

	return router.Run("localhost:8010")
}

func (ws *ApiStruct) Close() {

}
