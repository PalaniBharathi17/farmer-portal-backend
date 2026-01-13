package handlers

import "github.com/gin-gonic/gin"

// Health responds with service liveness.
func Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
