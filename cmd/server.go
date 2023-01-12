package main

import (
	"chasqi-go/cmd/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		handlePing(c)
	})
	r.Group("/gateway", gateway.HandleRequest)

	r.Run()
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
