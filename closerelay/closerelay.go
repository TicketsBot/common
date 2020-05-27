package closerelay

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type TicketClose struct {
	GuildId  uint64 `json:"guild_id"`
	TicketId int    `json:"ticket_id"`
	UserId   uint64 `json:"user_id"`
	Reason   string `json:"reason"`
}

const channel = "tickets:close"

func Publish(redis *redis.Client, data TicketClose) error {
	marshalled, err := json.Marshal(data); if err != nil {
		return err
	}
	return redis.Publish(channel, string(marshalled)).Err()
}

func Listen(redis *redis.Client, ch chan TicketClose) {
	for payload := range redis.Subscribe(channel).Channel() {
		var data TicketClose
		if err := json.Unmarshal([]byte(payload.Payload), &data); err != nil {
			continue
		}

		ch <- data
	}
}
