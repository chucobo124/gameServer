package rooms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gameserver/utils"
	cache "github.com/patrickmn/go-cache"

	"github.com/gin-gonic/gin"
)

// Room the room obj
type Room struct {
	Name         string `json:"name"`
	CurrentUsers int    `json:"current_user"`
	MaxUser      int    `json:"max_user"`
}

// RoomCache room cache
type RoomCache struct {
	ID    int `json:"id"`
	Users []User
}

// GetRoomsResp obj
type GetRoomsResp struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	LimitUsers int           `json:"limitUsers"`
	Creator    int           `json:"creator"`
	IsActive   bool          `json:"isActive"`
	CreDate    string        `json:"creDate"`
	UpdDate    string        `json:"updDate"`
	Game       []interface{} `json:"game"`
	User       []User        `json:"user"`
}

// User obj
type User struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Coin float32 `json:"coin"`
}

const externalRoomPath = "https://myfistwebsite-204102.appspot.com/api/Rooms/"
const externalUserPath = "https://myfistwebsite-204102.appspot.com/api/Users/"

var client = http.Client{}

// GetRooms Get Rooms Profile
func GetRooms(c *gin.Context) {
	userID := c.Query("user_id")

	rs, err := getRooms(c)
	if err != nil {
		utils.ReturnErrorJSON(c, err.Error())
		return
	}
	// Convert Rooms Object to response
	roomsCount := len(rs)
	rooms := make([]Room, roomsCount)
	for i, r := range rs {
		rooms[i].MaxUser = r.LimitUsers
		rooms[i].Name = r.Name
	}

	user, err := getUserProfile(c, userID)
	if err != nil {
		utils.ReturnErrorJSON(c, err.Error())
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
		utils.ReturnErrorJSON(c, err.Error())
		return
	} else if roomID == 0 {
		utils.ReturnErrorJSON(c, errors.New("The roomID is 0").Error())
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
		utils.ReturnErrorJSON(c, err.Error())
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
		utils.ReturnErrorJSON(c, err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ReturnErrorJSON(c, err.Error())
		return
	} else if resp.StatusCode != http.StatusOK {
		utils.ReturnErrorJSON(c, resp.Status)
		return
	}
	defer resp.Body.Close()

	c.JSON(200, gin.H{
		"message": "I want to join the room" + fmt.Sprint(roomID) + fmt.Sprint(
			reqForm),
	})
}

// JoinRoom save user into general cache
func JoinRoom(caches *cache.Cache) func(c *gin.Context) {
	return func(c *gin.Context) {
		roomID := c.Param("room_id")
		userID := c.Query("user_id")
		u, _ := strconv.Atoi(userID)
		if r, ok := caches.Get(roomID); ok {
			room := r.(*RoomCache)
			// check User is exist?
			var isUserExist bool
			for _, user := range room.Users {
				if user.ID == u {
					isUserExist = true
				}
			}
			if !isUserExist {
				user, err := getUserProfile(c, userID)
				if err != nil {
					utils.ReturnErrorJSON(c, err.Error())
					return
				}
				room.Users = append(room.Users, *user)
				caches.Set(roomID, room, cache.DefaultExpiration)
			}
		}
		c.JSON(200, nil)
	}
}

// GetCurrentRoom get current room profile
func GetCurrentRoom(caches *cache.Cache) func(c *gin.Context) {
	return func(c *gin.Context) {
		roomID := c.Param("room_id")
		if r, ok := caches.Get(roomID); ok {
			c.JSON(200, gin.H{
				"room": r,
			})
		} else {
			if room, err := getCurrentRoom(c, roomID); err != nil {
				utils.ReturnErrorJSON(c, err.Error())
				return
			} else {
				rc := new(RoomCache)
				rc.ID, _ = strconv.Atoi(roomID)
				for _, user := range room.User {
					rc.Users = append(rc.Users, user)
				}
				caches.Set(roomID, rc, cache.DefaultExpiration)
				if room, ok := caches.Get(roomID); ok {
					c.JSON(200, gin.H{
						"room": room,
					})
				} else {
					utils.ReturnErrorJSON(c, "Cache Dose Not save successful")
				}
			}
		}
	}
}

func getCurrentRoom(c *gin.Context, roomID string) (*GetRoomsResp, error) {
	req, err := http.NewRequest(http.MethodGet, externalRoomPath+roomID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	getRoomResp := new(GetRoomsResp)
	if err := json.Unmarshal(body, getRoomResp); err != nil {
		return nil, err
	}
	return getRoomResp, nil
}

func getRooms(c *gin.Context) ([]GetRoomsResp, error) {
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
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	getRoomResp := make([]GetRoomsResp, 0)
	if err := json.Unmarshal(body, &getRoomResp); err != nil {
		return nil, err
	}
	return getRoomResp, nil
}

func getUserProfile(c *gin.Context, userID string) (*User, error) {

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
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	currentUser := new(User)
	if err := json.Unmarshal(body, currentUser); err != nil {
		return nil, err
	}
	return currentUser, nil

}
