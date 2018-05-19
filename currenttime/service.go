package currenttime

import (
	"time"

	"github.com/gin-gonic/gin"
)

// GetTime Get current time as time server
func GetTime(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": time.Now().Unix(),
	})
}
