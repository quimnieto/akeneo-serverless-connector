package handlers

import (
	"akeneo-event-config/internal/akeneo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	client *akeneo.Client
}

func NewHandler(client *akeneo.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) GetSubscriber(c *gin.Context) {
	subscriber, err := h.client.GetSubscriber()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subscriber)
}

func (h *Handler) CreateSubscriber(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.client.CreateSubscriber(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Subscriber created"})
}

func (h *Handler) UpdateSubscriber(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.client.UpdateSubscriber(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscriber updated"})
}

func (h *Handler) GetSubscriptions(c *gin.Context) {
	subscriptions, err := h.client.GetSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subscriptions)
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.client.CreateSubscription(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Subscription created"})
}

func (h *Handler) UpdateSubscription(c *gin.Context) {
	connectionCode := c.Param("code")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.client.UpdateSubscription(connectionCode, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated"})
}

func (h *Handler) DeleteSubscription(c *gin.Context) {
	connectionCode := c.Param("code")
	if err := h.client.DeleteSubscription(connectionCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
}

func (h *Handler) GetEventTypes(c *gin.Context) {
	eventTypes := h.client.GetEventTypes()
	c.JSON(http.StatusOK, eventTypes)
}
