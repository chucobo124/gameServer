package rooms

import (
	"github.com/gin-gonic/gin"
)

// GetRooms Get Rooms Profile
func GetRooms(c *gin.Context) {
	type room struct {
		Name         string `json:"name"`
		CurrentUsers int    `json:"current_user"`
		MaxUser      int    `json:"max_user"`
	}
	type user struct {
		Name  string `json:"name"`
		Coins int    `json:"coins"`
	}
	r := []room{
		room{
			Name:         "RoomA",
			CurrentUsers: 10,
			MaxUser:      20,
		},
		room{
			Name:         "RoomB",
			CurrentUsers: 10,
			MaxUser:      20,
		},
		room{
			Name:         "RoomC",
			CurrentUsers: 10,
			MaxUser:      20,
		},
		room{
			Name:         "RoomA",
			CurrentUsers: 10,
			MaxUser:      20,
		},
	}
	u := user{
		Name:  "Lucas",
		Coins: 10,
	}
	c.JSON(200, gin.H{
		"rooms": r,
		"user":  u,
	})
}

// PutRoom User can join the room
func PutRoom(c *gin.Context) {
	// roomId := c.Param("room_id")
	// Check room is exist?

	c.JSON(200, gin.H{
		"message": "I want to join the room" + c.Param("room_id"),
	})
}
