package chatrelay

import (
	"encoding/json"
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis"
	"github.com/rxdn/gdl/objects/channel/message"
)

type MessageData struct {
	Ticket  database.Ticket `json:"ticket"`
	Message message.Message `json:"message"`
}

const channel = "tickets:chatrelay"

func PublishMessage(redis *redis.Client, data MessageData) error {
	marshalled, err := json.Marshal(data); if err != nil {
		return err
	}

	return redis.Publish(channel, string(marshalled)).Err()
}

func Listen(redis *redis.Client, ch chan MessageData) {
	for payload := range redis.Subscribe(channel).Channel() {
		var data MessageData

		if err := json.Unmarshal([]byte(payload.Payload), &data); err != nil {
			continue
		}

		ch <- data
	}
}
