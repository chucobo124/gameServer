package rooms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/gin-gonic/gin"
)

// type getRoomsResp struct {
// 	ID         int           `json:"id"`
// 	Name       string        `json:"name"`
// 	LimitUsers int           `json:"limitUsers"`
// 	Creator    string        `json:"creator"`
// 	IsActive   bool          `json:"isActive"`
// 	CreDate    string        `json:"creDate"`
// 	UpdDate    string        `json:"updDate"`
// 	Game       []interface{} `json:"game"`
// 	User       []user        `json:"user"`
// }

type GetRoomsResp struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	LimitUsers int           `json:"limitUsers"`
	Creator    string        `json:"creator"`
	IsActive   bool          `json:"isActive"`
	CreDate    string        `json:"creDate"`
	UpdDate    string        `json:"updDate"`
	Game       []interface{} `json:"game"`
	User       []interface{} `json:"user"`
}

type User struct {
	Name  string `json:"name"`
	Coins int    `json:"coins"`
}

const externalRoomPath = "https://myfistwebsite-204102.appspot.com/api/Rooms/"
const externalUserPath = "https://myfistwebsite-204102.appspot.com/api/Users/"

// GetRooms Get Rooms Profile
func GetRooms(c *gin.Context) {
	userID := c.Query("user_id")
	type room struct {
		Name         string `json:"name"`
		CurrentUsers int    `json:"current_user"`
		MaxUser      int    `json:"max_user"`
	}

	rs, err := getRooms(c)
	if err != nil {
		returnErrorJSON(c, err.Error())
		return
	}
	// Convert Rooms Object to response
	roomsCount := len(rs)
	rooms := make([]room, roomsCount)
	for i, r := range rs {
		rooms[i].MaxUser = r.LimitUsers
		rooms[i].Name = r.Name
	}

	user, err := getUserProfile(c, userID)
	if err != nil {
		returnErrorJSON(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"rooms": rooms,
		"user":  user,
	})
}

// PutRoom User can join the room
func PutRoom(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("room_id"))
	if err != nil {
		returnErrorJSON(c, err.Error())
		return
	} else if roomID == 0 {
		returnErrorJSON(c, errors.New("The roomID is 0").Error())
		return
	}

	putRoomPath := externalRoomPath + fmt.Sprint(roomID)

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
		return
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
	debug.PrintStack()
	c.JSON(200, gin.H{
		"error": msg,
	})
}

func getRooms(c *gin.Context) ([]GetRoomsResp, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, externalRoomPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	getRoomResp := make([]GetRoomsResp, 0)
	if err := json.Unmarshal(body, &getRoomResp); err != nil {
		return nil, err
	}
	return getRoomResp, nil
}

func getUserProfile(c *gin.Context, userID string) (*User, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, externalUserPath+userID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	currentUser := new(User)
	if err := json.Unmarshal(body, &currentUser); err != nil {
		return nil, err
	}
	return currentUser, nil

}
