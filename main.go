// Provide time and room for client
package main

import (
	"github.com/gameserver/currenttime"
	"github.com/gameserver/rooms"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/time", currenttime.GetTime)
	r.GET("/rooms", rooms.GetRooms)
	r.PUT("/room/:room_id", rooms.PutRoom)
	r.Run() // listen and serve on 0.0.0.0:8080
}
