package utils

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func ReturnErrorJSON(c *gin.Context, msg string) {
	debug.PrintStack()
	c.JSON(200, gin.H{
		"error": msg,
	})
}
