package service

import (
	"server_order/structure"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type service struct {
	producer sarama.SyncProducer
}

func New(producer sarama.SyncProducer) *service {
	return &service{producer}
}

func (s *service) Buy(c *gin.Context) {
	var req structure.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "input error"})
		return
	}

	if req.Type == "sell" {
		c.AbortWithStatusJSON(400, gin.H{"error": "action error"})
		return
	}

	if err := s.handleOrder(&req); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "server error"})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (s *service) Sell(c *gin.Context) {
	var req structure.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "input error"})
		return
	}

	if req.Type == "buy" {
		c.AbortWithStatusJSON(400, gin.H{"error": "action error"})
		return
	}

	if err := s.handleOrder(&req); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "server error"})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (s *service) Test(c *gin.Context) {
	if err := s.handleTest(); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "server error"})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}
