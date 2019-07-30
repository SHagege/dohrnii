package handlers

import (
	"dohrnii/internal/app/block"
	"github.com/gin-gonic/gin"
)

// GetBlockchain ...
func GetBlockchain(c *gin.Context) {
	c.JSON(200, gin.H{
		"blockchain": block.Bc.Chain,
	})
}

func GetLastBlock(c *gin.Context) {
	c.JSON(200, gin.H{
		"last_block": block.Bc.GetLastBlock(),
	})
}