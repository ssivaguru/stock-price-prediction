package apiserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ApiInterface interface {
	StartServer() error
	Close()
}

type ApiStruct struct {
	isClose chan bool
}

func New() ApiInterface {
	return &ApiStruct{isClose: make(chan bool)}
}

func handlePredit(c *gin.Context, ws *ApiStruct) {
	paramPairs := c.Request.URL.Query()
	for key, values := range paramPairs {
		fmt.Printf("key = %v, value(s) = %v\n", key, values)
	}

	c.JSON(200, "{Status:Success}")
}

func handleDescribe(c *gin.Context, ws *ApiStruct) {
	paramPairs := c.Request.URL.Query()
	for key, values := range paramPairs {
		fmt.Printf("key = %v, value(s) = %v\n", key, values)
	}
	c.JSON(200, "{Status:Success}")
}

func (ws *ApiStruct) setupRoutes(router *gin.Engine) {
	router.GET("/api/predict", func(c *gin.Context) {
		handlePredit(c, ws)
	})

	router.GET("/api/describe", func(c *gin.Context) {
		handleDescribe(c, ws)
	})
}

func (ws *ApiStruct) StartServer() error {
	router := gin.Default()
	ws.setupRoutes(router)
	router.Run("localhost:8010")

	return nil
}

func (ws *ApiStruct) Close() {

}
