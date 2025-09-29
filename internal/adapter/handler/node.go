package handler

import (
	"kubecompute/internal/core/port"

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	service port.NodeService
}

func NewNodeHandler(service port.NodeService) *NodeHandler {
	return &NodeHandler{
		service: service,
	}
}

func (h *NodeHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/nodes", h.ListNodes)
	router.POST("/nodes", h.CreateNode)
	router.GET("/nodes/:name", h.GetNode)
	router.PUT("/nodes/:name", h.UpdateNode)
	router.DELETE("/nodes/:name", h.DeleteNode)
}

func (h *NodeHandler) ListNodes(c *gin.Context) {
	// Implementation omitted for brevity
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
	// Implementation omitted for brevity
}

func (h *NodeHandler) GetNode(c *gin.Context) {
	// Implementation omitted for brevity
}

func (h *NodeHandler) UpdateNode(c *gin.Context) {
	// Implementation omitted for brevity
}

func (h *NodeHandler) DeleteNode(c *gin.Context) {
	// Implementation omitted for brevity
}
