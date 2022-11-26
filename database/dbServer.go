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

func (s *DbServer) handlePredict(c *gin.Context) {

	//we need to do more handling
}

func (s *DbServer) handleCreateStock(c *gin.Context) {
	var bodyMap map[string]string
	json.NewDecoder(c.Request.Body).Decode(&bodyMap)

	if _, ok := bodyMap["name"]; !ok {
		c.JSON(400, "{Msg: Malformed URL}")
	}

	//jsut check if its already created

	_, err := s.Dbclient.GetData(bodyMap["name"])

	if err == nil {
		//then stock is alreadty created
		c.JSON(501, "{Msg:Stock Is already created}")
	}

	err = s.Dbclient.InsertData(&database.Stock{ID: primitive.NewObjectIDFromTimestamp(time.Now()), Name: bodyMap["name"], Path: "", Status: StatusTraning, UpdatedAt: time.Now()})

	if err != nil {
		log.Println("handleCreateStock:: error creating stock ", err)
		c.JSON(500, "{Msg:Error creating stock}")
	}

}

func (s *DbServer) handleDescribe(c *gin.Context) {
}

func (s *DbServer) handleDeleteStock(c *gin.Context) {
}

func (s *DbServer) handleUpdateStock(c *gin.Context) {
}

func (s *DbServer) setupRoutes(router *gin.Engine) {
	router.GET("/api/predict", func(c *gin.Context) {
		s.handlePredict(c)
	})

	router.GET("/api/describe", func(c *gin.Context) {
		s.handleDescribe(c)
	})

	router.POST("/api/stock", func(c *gin.Context) {
		s.handleCreateStock(c)
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
