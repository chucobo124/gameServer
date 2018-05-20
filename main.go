// Provide time and room for client
package main

import (
	"time"

	"github.com/gameserver/currenttime"
	"github.com/gameserver/rooms"
	"github.com/gin-gonic/gin"
	cache "github.com/patrickmn/go-cache"
)

func main() {
	c := cache.New(5*time.Minute, 10*time.Minute)
	r := gin.Default()
	r.GET("/join/:room_id", rooms.JoinRoom(c))
	r.GET("/room/:room_id", rooms.GetCurrentRoom(c))
	r.GET("/time", currenttime.GetTime)
	r.GET("/rooms", rooms.GetRooms)
	r.PUT("/room/:room_id", rooms.PutRoom)
	r.Run() // listen and serve on 0.0.0.0:8080
}
