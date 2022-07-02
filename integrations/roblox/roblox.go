package roblox

import (
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/common/webproxy"
	"net/http"
)

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	IsBanned    bool   `json:"isBanned"`
	Created     string `json:"created"`
}

func RequestUserData(proxy *webproxy.WebProxy, robloxId int) (User, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://users.roblox.com/v1/users/%d", robloxId), nil)
	if err != nil {
		return User{}, err
	}

	res, err := proxy.Do(req)
	if err != nil {
		return User{}, err
	}

	var data User
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return User{}, err
	}

	return data, nil
}
