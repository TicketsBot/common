package bloxlink

import (
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/common/webproxy"
	"net/http"
)

type BloxlinkResponse struct {
	Success bool `json:"success"`
	User    struct {
		RobloxId       int `json:"robloxId,string"`
		PrimaryAccount int `json:"primaryAccount,string"`
	} `json:"user"`
	Error string `json:"error"`
}

var (
	ErrQuotaExceeded = fmt.Errorf("Bloxlink API quota exceeded")
	ErrUserNotFound  = fmt.Errorf("User not found")
)

func RequestUserId(proxy *webproxy.WebProxy, bloxlinkApiKey string, userId uint64) (int, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://v3.blox.link/developer/discord/%d", userId), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("api-key", bloxlinkApiKey)

	res, err := proxy.Do(req)
	if err != nil {
		return 0, err
	}

	var response BloxlinkResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return 0, err
	}

	if !response.Success {
		if res.Header.Get("Quota-Remaining") == "0" {
			return 0, ErrQuotaExceeded
		} else if response.Error == "The specified user is not verified with Bloxlink." {
			return 0, ErrUserNotFound
		} else {
			return 0, fmt.Errorf("Bloxlink API request unsuccessful: %s", response.Error)
		}
	}

	if response.User.RobloxId == 0 {
		return 0, ErrUserNotFound
	}

	return response.User.RobloxId, nil
}
