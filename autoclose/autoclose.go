package chatrelay

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type Ticket struct {
	GuildId  uint64
	TicketId int
}

const channel = "tickets:autoclose"

func PublishMessage(redis *redis.Client, data []Ticket) error {
	marshalled, err := json.Marshal(data); if err != nil {
		return err
	}

	return redis.Publish(channel, string(marshalled)).Err()
}

func Listen(redis *redis.Client, ch chan Ticket) {
	for payload := range redis.Subscribe(channel).Channel() {
		var data Ticket

		if err := json.Unmarshal([]byte(payload.Payload), &data); err != nil {
			continue
		}

		ch <- data
	}
}
