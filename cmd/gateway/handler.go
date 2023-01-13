package gateway

import (
	"chasqi-go/core/engine"
	"chasqi-go/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	coreEngine *engine.DefaultEngine
}

func NewHandler(coreEngine *engine.DefaultEngine) *Handler {
	return &Handler{coreEngine: coreEngine}
}

func (h *Handler) Handle(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		h.Get(c)
	case http.MethodPost:
		h.Post(c)
	}
}

func (h *Handler) Get(c *gin.Context) {
	log.Printf("Path: %s", c.Request.URL.Path)
	if strings.Contains(c.Request.URL.Path, "status") {
		treeId := c.Param("treeId")
		status := h.coreEngine.LoopStatus(treeId)

		if status == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}

		c.JSON(http.StatusOK, status)
	}
	if strings.Contains(c.Request.URL.Path, "result") {
		log.Printf("result")
		treeId := c.Param("treeId")
		result := h.coreEngine.TestResult(treeId)
		if result == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) Post(c *gin.Context) {
	treeRequest := &types.Tree{}

	if err := c.BindJSON(treeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	treeRequest.ID = uuid.NewString()
	if err := h.coreEngine.Enqueue(treeRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, treeRequest)
}
