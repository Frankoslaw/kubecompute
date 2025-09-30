package handler

import (
	"kubecompute/internal/core/domain"
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
	router.GET("/nodes", h.ListNodesAll)
	router.GET("/ns/:namespace/nodes", h.ListNodes)
	router.POST("/ns/:namespace/nodes", h.CreateNode)
	router.GET("/ns/:namespace/nodes/:name", h.GetNode)
	router.PUT("/ns/:namespace/nodes/:name", h.UpdateNode)
	router.DELETE("/ns/:namespace/nodes/:name", h.DeleteNode)
}

// ListNodesAll godoc
//
//	@Summary		List all nodes across all namespaces
//	@Description	Returns all nodes, ignoring namespace
//	@Tags			nodes
//	@Produce		json
//	@Success		200	{array}		domain.Node
//	@Failure		500	{object}	map[string]string
//	@Router			/nodes [get]
func (h *NodeHandler) ListNodesAll(c *gin.Context) {
	nodes, err := h.service.ListNodes(c.Request.Context(), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, nodes)
}

// ListNodes godoc
//
//	@Summary		List nodes in a namespace
//	@Description	Returns nodes filtered by namespace
//	@Tags			nodes
//	@Produce		json
//	@Param			namespace	path		string	true	"Namespace"
//	@Success		200			{array}		domain.Node
//	@Failure		500			{object}	map[string]string
//	@Router			/ns/{namespace}/nodes [get]
func (h *NodeHandler) ListNodes(c *gin.Context) {
	namespace := c.Param("namespace")
	nodes, err := h.service.ListNodes(c.Request.Context(), &namespace)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, nodes)
}

// CreateNode godoc
//
//	@Summary		Create a node
//	@Description	Create a new node in the given namespace
//	@Tags			nodes
//	@Accept			json
//	@Produce		json
//	@Param			namespace	path		string		true	"Namespace"
//	@Param			node		body		domain.Node	true	"Node object"
//	@Success		201			{object}	domain.Node
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/ns/{namespace}/nodes [post]
func (h *NodeHandler) CreateNode(c *gin.Context) {
	var node *domain.Node = &domain.Node{}
	if err := c.ShouldBindJSON(node); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	node, err := h.service.CreateNode(c.Request.Context(), node)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, node)
}

// GetNode godoc
//
//	@Summary		Get a node
//	@Description	Returns a single node by namespace and name
//	@Tags			nodes
//	@Produce		json
//	@Param			namespace	path		string	true	"Namespace"
//	@Param			name		path		string	true	"Node name"
//	@Success		200			{object}	domain.Node
//	@Failure		500			{object}	map[string]string
//	@Router			/ns/{namespace}/nodes/{name} [get]
func (h *NodeHandler) GetNode(c *gin.Context) {
	name := domain.NamespacedName{
		Namespace: c.Param("namespace"),
		Name:      c.Param("name"),
	}
	node, err := h.service.GetNode(c.Request.Context(), name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if node == nil {
		c.JSON(404, gin.H{"error": "node not found"})
		return
	}
	c.JSON(200, node)
}

// UpdateNode godoc
//
//	@Summary		Update a node
//	@Description	Updates a node in a namespace
//	@Tags			nodes
//	@Accept			json
//	@Produce		json
//	@Param			namespace	path		string		true	"Namespace"
//	@Param			name		path		string		true	"Node name"
//	@Param			node		body		domain.Node	true	"Node object"
//	@Success		200			{object}	domain.Node
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/ns/{namespace}/nodes/{name} [put]
func (h *NodeHandler) UpdateNode(c *gin.Context) {
	var node *domain.Node = &domain.Node{}
	if err := c.ShouldBindJSON(node); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	node.ObjectMeta.Namespace = c.Param("namespace")
	node.ObjectMeta.Name = c.Param("name")

	node, err := h.service.UpdateNode(c.Request.Context(), node)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, node)
}

// DeleteNode godoc
//
//	@Summary		Delete a node
//	@Description	Soft deletes a node in a namespace
//	@Tags			nodes
//	@Produce		json
//	@Param			namespace	path	string		true	"Namespace"
//	@Param			name		path	string		true	"Node name"
//	@Param			node		body	domain.Node	true	"Node object"
//	@Success		204			"No Content"
//	@Failure		500			{object}	map[string]string
//	@Router			/ns/{namespace}/nodes/{name} [delete]
func (h *NodeHandler) DeleteNode(c *gin.Context) {
	var node *domain.Node = &domain.Node{}
	if err := c.ShouldBindJSON(node); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	node.ObjectMeta.Namespace = c.Param("namespace")
	node.ObjectMeta.Name = c.Param("name")

	_, err := h.service.DeleteNode(c.Request.Context(), node)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}
