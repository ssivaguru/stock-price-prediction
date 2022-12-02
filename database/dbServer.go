package main

import (
	"encoding/json"
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

	if _, ok := bodyMap["name"]; !ok {
		c.JSON(400, "{Msg: Malformed URL}")
		return
	}

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

	err = s.Dbclient.InsertData(&database.Stock{ID: primitive.NewObjectIDFromTimestamp(time.Now()), Name: bodyMap.Get("name"), Path: "", Status: StatusTraning, UpdatedAt: time.Now()})

	if err != nil {
		log.Println("handleCreateStock:: error creating stock ", err)
		c.Data(500, gin.MIMEJSON, getJson(map[string]string{"error": "error creating stock record"}))
		return
	}

	c.Data(201, gin.MIMEJSON, getJson(map[string]string{"error": "Created stock record"}))
}

func (s *DbServer) handleDescribe(c *gin.Context) {
}

func (s *DbServer) handleDeleteStock(c *gin.Context) {
}

func (s *DbServer) handleUpdateStock(c *gin.Context) {
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
	s.setupRoutes(router)
	return router.Run("localhost:8080")
}
