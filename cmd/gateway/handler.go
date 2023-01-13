package gateway

import (
	"chasqi-go/core/engine"
	"chasqi-go/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Handler struct {
	coreEngine *engine.DefaultEngine
}

func NewHandler(coreEngine *engine.DefaultEngine) *Handler {
	return &Handler{coreEngine: coreEngine}
}

func (h *Handler) Get(c *gin.Context) {
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
