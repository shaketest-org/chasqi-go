package gateway

import (
	"chasqi-go/core/engine"
	"chasqi-go/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Handler struct {
	coreEngine *engine.DefaultEngine
}

func NewHandler(coreEngine *engine.DefaultEngine) *Handler {
	return &Handler{coreEngine: coreEngine}
}

func (h *Handler) HandleRequest(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		h.handleGet(c)
	case "POST":
		h.handlePost(c)
	}
}

func (h *Handler) handleGet(c *gin.Context) {
	if strings.Contains("status", c.Request.URL.Path) {
		treeId := c.Param("treeId")
		status := h.coreEngine.ById(treeId)
		if status == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}

		c.JSON(http.StatusOK, status)
	}
	if strings.Contains("result", c.Request.URL.Path) {
		treeId := c.Param("treeId")
		result := h.coreEngine.Get(treeId)
		if result == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) handlePost(c *gin.Context) {
	var treeRequest types.Tree
	if err := c.ShouldBindJSON(&treeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.coreEngine.Enqueue(&treeRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "accepted"})
}
