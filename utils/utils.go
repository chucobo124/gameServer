package utils

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ReturnErrorJSON Util Error msg function
func ReturnErrorJSON(c *gin.Context, msg string) {
	fmt.Print(msg)
	debug.PrintStack()
	c.JSON(200, gin.H{
		"error": msg,
	})
}
