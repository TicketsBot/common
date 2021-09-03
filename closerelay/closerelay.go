package closerelay

import (
	"context"
	"encoding/json"
	"github.com/TicketsBot/common/utils"
	"github.com/go-redis/redis/v8"
)

type TicketClose struct {
	GuildId        uint64 `json:"guild_id"`
	TicketId       int    `json:"ticket_id"`
	UserId         uint64 `json:"user_id"`
	Reason         string `json:"reason"`
}

const key = "tickets:close"

func Publish(redis *redis.Client, data TicketClose) error {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return redis.RPush(utils.DefaultContext(), key, string(marshalled)).Err()
}

func Listen(redis *redis.Client, ch chan TicketClose) {
	for {
		res, err := redis.BLPop(context.Background(), 0, key).Result()
		if err != nil {
			continue
		}

		var data TicketClose
		if err := json.Unmarshal([]byte(res[1]), &data); err != nil {
			continue
		}

		ch <- data
	}
}
