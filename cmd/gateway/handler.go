package gateway

import "github.com/gin-gonic/gin"

func HandleRequest(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		handleGet(c)
	}
}

func handleGet(c *gin.Context) {

}
