package rooms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	roomID, err := strconv.Atoi(c.Param("room_id"))
	if err != nil {
		returnErrorJSON(c, err.Error())
	} else if roomID == 0 {
		returnErrorJSON(c, errors.New("The roomID is 0").Error())
	}

	putRoomPath := "https://myfistwebsite-204102.appspot.com/api/Room/" + fmt.Sprint(roomID)

	type requestForm struct {
		Name     string `json:"room_name" `
		MaxUsers int    `json:"limit_users"`
		Creator  string `json:"creator"`
		IsActive bool   `json:"is_active"`
	}
	reqForm := new(requestForm)

	type putRoomRequest struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		MaxUsers int    `json:"limitUsers"`
		Creator  string `json:"creator"`
		IsActive bool   `json:"isActive"`
	}

	err = c.BindJSON(reqForm)
	if err != nil {
		returnErrorJSON(c, err.Error())
	}

	// Check validation
	putRoomReq := putRoomRequest{}
	if !reqForm.IsActive {
		putRoomReq.IsActive = reqForm.IsActive
	}

	if reqForm.Creator != "" {
		putRoomReq.Creator = reqForm.Creator
	}

	if reqForm.MaxUsers != 0 {
		putRoomReq.MaxUsers = reqForm.MaxUsers
	}

	if reqForm.Name != "" {
		putRoomReq.Name = reqForm.Name
	}

	jsonPutRoomReq, _ := json.Marshal(putRoomReq)
	// Check room is exist?
	req, err := http.NewRequest(http.MethodPost, putRoomPath, bytes.NewBuffer(jsonPutRoomReq))
	if err != nil {
		returnErrorJSON(c, err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		returnErrorJSON(c, err.Error())
		return
	} else if resp.StatusCode != http.StatusOK {
		returnErrorJSON(c, resp.Status)
		return
	}

	c.JSON(200, gin.H{
		"message": "I want to join the room" + fmt.Sprint(roomID) + fmt.Sprint(
			reqForm),
	})
}

func returnErrorJSON(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"error": msg,
	})
}
