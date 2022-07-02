package roblox

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/common/webproxy"
	"github.com/go-redis/redis/v8"
	"time"
)

type BloxlinkIntegration struct {
	redis  *redis.Client
	proxy  *webproxy.WebProxy
	apiKey string
}

func NewBloxlinkIntegration(redis *redis.Client, proxy *webproxy.WebProxy, apiKey string) *BloxlinkIntegration {
	return &BloxlinkIntegration{
		redis:  redis,
		proxy:  proxy,
		apiKey: apiKey,
	}
}

func (i *BloxlinkIntegration) GetRobloxUser(discordUserId uint64) (User, error) {
	redisKey := fmt.Sprintf("bloxlink:%d", discordUserId)

	cached, err := i.redis.Get(context.Background(), redisKey).Result()
	if err == nil {
		var user User
		if err := json.Unmarshal([]byte(cached), &user); err != nil {
			return User{}, err
		}

		return user, nil
	} else if err != redis.Nil { // If the error is redis.Nil, this means that the key does not exist, and we should continue
		return User{}, err
	}

	robloxId, err := RequestUserId(i.proxy, i.apiKey, discordUserId)
	if err != nil {
		return User{}, err
	}

	user, err := RequestUserData(i.proxy, robloxId)
	if err != nil {
		return User{}, err
	}

	go func() {
		encoded, err := json.Marshal(user)
		if err != nil {
			return
		}

		i.redis.SetEX(context.Background(), redisKey, string(encoded), time.Hour*24)
	}()

	return user, nil
}
