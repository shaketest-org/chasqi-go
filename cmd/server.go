package main

import (
	"chasqi-go/cmd/gateway"
	"chasqi-go/core/agent"
	"chasqi-go/core/engine"
	"chasqi-go/visitor"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	e := engine.New(
		func() agent.NodeVisitor { return visitor.NewDefaultHttpClient() },
		make(chan struct{}),
	)
	h := gateway.NewHandler(e)

	r.GET("/ping", func(c *gin.Context) {
		handlePing(c)
	})

	r.Group("/gateway", h.HandleRequest)

	go func() {
		e.Start()
	}()

	r.Run()
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
