package bloxlink

import (
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/common/webproxy"
	"github.com/pkg/errors"
	"net/http"
)

type BloxlinkResponse struct {
	RobloxId int `json:"robloxID,string"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrQuotaExceeded = fmt.Errorf("Bloxlink API quota exceeded")
	ErrUserNotFound  = fmt.Errorf("User not found")
)

func RequestUserId(proxy *webproxy.WebProxy, bloxlinkApiKey string, userId uint64) (int, error) {
	url := fmt.Sprintf("https://api.blox.link/v4/public/discord-to-roblox/%d", userId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", bloxlinkApiKey)

	res, err := proxy.Do(req)
	if err != nil {
		return 0, err
	}

	switch res.StatusCode {
	case http.StatusOK:
		break // continue
	case http.StatusNotFound:
		return 0, ErrUserNotFound
	case http.StatusTooManyRequests:
		return 0, ErrQuotaExceeded
	default:
		var errorResponse ErrorResponse

		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			return 0, errors.Wrapf(err, "failed to decode bloxlink error response - status code was %d", res.StatusCode)
		}

		return 0, errors.Wrap(errors.New(errorResponse.Error), "bloxlink api returned error")
	}

	var response BloxlinkResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, err
	}

	if response.RobloxId == 0 {
		return 0, ErrUserNotFound
	}

	return response.RobloxId, nil
}
