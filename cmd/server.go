package main

import (
	"chasqi-go/cmd/gateway"
	"chasqi-go/core/agent"
	"chasqi-go/core/engine"
	"chasqi-go/data"
	"chasqi-go/data/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	e := engine.New(
		func() agent.NodeVisitor { return data.NewDefaultHttpClient() },
		result.NewManager(),
		make(chan struct{}),
	)
	h := gateway.NewHandler(e)

	r.GET("/ping", func(c *gin.Context) {
		handlePing(c)
	})

	gateway := r.Group("/gateway")
	{
		gateway.GET("", h.Get)
		gateway.POST("", h.Post)
	}

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
