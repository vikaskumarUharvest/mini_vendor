package handler

import (
	"github.com/gin-gonic/gin"

	"pgxpostgress/domain"
	"pgxpostgress/service"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{svc: s}
}
func (h *OrderHandler) Create(c *gin.Context) {
	var o domain.Order

	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	err := h.svc.Create(c, &o)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, o)
}

func (h *OrderHandler) Get(c *gin.Context) {
	id := c.Param("id")

	order, err := h.svc.Get(c, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(200, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	orders, _ := h.svc.List(c)
	c.JSON(200, orders)
}

func (h *OrderHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Status string `json:"status"`
	}

	c.ShouldBindJSON(&body)

	h.svc.Update(c, id, body.Status)

	c.JSON(200, gin.H{"updated": true})
}

func (h *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	h.svc.Delete(c, id)

	c.JSON(200, gin.H{"deleted": true})
}