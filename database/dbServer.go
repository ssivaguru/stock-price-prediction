package main

import (
	"encoding/json"
	"fmt"
	"log"
	"stock-price-prediction/database/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DbServer struct {
	Dbclient database.DbInterface
}

func getJson(data map[string]string) []byte {
	jsonData, _ := json.Marshal(data)
	log.Println(string(jsonData[:]))
	return jsonData
}
func (s *DbServer) handlePredict(c *gin.Context) {

	bodyMap := c.Request.URL.Query()
	log.Println(bodyMap)

	if _, ok := bodyMap["name"]; !ok {
		c.JSON(400, "{Msg: Malformed URL}")
		return
	}

	fmt.Println(" predicting for data ", bodyMap.Get("name"))
	//jsut check if its already created
	data, err := s.Dbclient.GetData(bodyMap.Get("name"))

	if err == nil {
		outByte, err := json.Marshal(data)
		if err != nil {
			log.Println("error encoding jdon ", err)
			c.Data(400, gin.MIMEJSON, getJson(map[string]string{"error": "error encoding json"}))
			return
		}
		log.Println(string(outByte[:]))
		c.Data(200, gin.MIMEJSON, outByte)
		return
	}

	newStockData := &database.Stock{ID: primitive.NewObjectIDFromTimestamp(time.Now()), Name: bodyMap.Get("name"), Path: "", Status: StatusTraning, UpdatedAt: time.Now(), Requesting_user: nil, Predicted_val: nil}
	err = s.Dbclient.InsertData(newStockData)

	if err != nil {
		log.Println("handleCreateStock:: error creating stock ", err)
		c.Data(500, gin.MIMEJSON, getJson(map[string]string{"error": "error creating stock record"}))
		return
	}

	outByte, err := json.Marshal(newStockData)
	c.Data(200, gin.MIMEJSON, outByte)
}

func (s *DbServer) handleDescribe(c *gin.Context) {
}

func (s *DbServer) handleDeleteStock(c *gin.Context) {
}

func (s *DbServer) handleUpdateStock(c *gin.Context) {
	var buffer map[string]interface{}

	json.NewDecoder(c.Request.Body).Decode(&buffer)

	log.Println(buffer)

	respos := buffer["predict"].([]interface{})
	var prediction []float64

	for _, val := range respos {
		prediction = append(prediction, val.(float64))
	}
	fmt.Println(prediction)

	s.Dbclient.UpdateByName(buffer["name"].(string), prediction)
	s.Dbclient.UpdateState(buffer["name"].(string), StatusReady)

	c.Data(200, gin.MIMEJSON, getJson(map[string]string{"status": " updating value"}))
}

func (s *DbServer) setupRoutes(router *gin.Engine) {
	//will handle getting and creating data
	router.GET("/api/predict", func(c *gin.Context) {
		s.handlePredict(c)
	})

	router.GET("/api/describe", func(c *gin.Context) {
		s.handleDescribe(c)
	})

	router.DELETE("/api/stock", func(c *gin.Context) {
		s.handleDeleteStock(c)
	})

	router.PUT("/api/stock", func(c *gin.Context) {
		s.handleUpdateStock(c)
	})
}

func (s *DbServer) start() error {

	router := gin.Default()
	fmt.Println("setting default route")
	s.setupRoutes(router)
	fmt.Println("starting router")
	return router.Run("localhost:8080")
}
