package main

import (
	"chasqi-go/cmd/gateway"
	"chasqi-go/core/agent"
	"chasqi-go/core/engine"
	"chasqi-go/data"
	"chasqi-go/data/result"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
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

	gw := r.Group("/gateway")
	{
		gw.GET("/status/:treeId", h.Get)
		gw.GET("/result/:treeId", h.Get)
		gw.POST("", h.Post)
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
