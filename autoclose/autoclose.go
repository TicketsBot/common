package autoclose

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type Ticket struct {
	GuildId  uint64 `json:"guild_id"`
	TicketId int    `json:"ticket_id"`
}

const channel = "tickets:autoclose"

func PublishMessage(redis *redis.Client, data []Ticket) error {
	var marshalled []interface{}
	for _, ticket := range data {
		json, err := json.Marshal(ticket)
		if err != nil {
			return err
		}

		marshalled = append(marshalled, string(json))
	}

	return redis.RPush(channel, marshalled...).Err()
}

func Listen(redis *redis.Client, ch chan Ticket) {
	for {
		data, err := redis.BLPop(0, channel).Result()
		if err != nil || len(data) < 2 {
			continue
		}

		// data = [list_name, content]
		var ticket Ticket
		if err := json.Unmarshal([]byte(data[1]), &ticket); err != nil {
			continue
		}

		ch <- ticket
	}
}
