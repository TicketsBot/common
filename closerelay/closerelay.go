package closerelay

import (
	"encoding/json"
	"github.com/go-redis/redis"
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
	return redis.RPush(key, string(marshalled)).Err()
}

func Listen(redis *redis.Client, ch chan TicketClose) {
	for {
		res, err := redis.BLPop(0, key).Result()
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
