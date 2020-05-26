package tokenchange

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

const channel = "tickets:tokenchange"

type TokenChangeData struct {
	Token string `json:"token"`
	NewId uint64 `json:"new_id"`
	OldId uint64 `json:"old_id"`
}

func PublishTokenChange(client *redis.Client, data TokenChangeData) error {
	marshalled, err := json.Marshal(data); if err != nil {
		return err
	}

	return client.Publish(channel, string(marshalled)).Err()
}

func ListenTokenChange(client *redis.Client, ch chan TokenChangeData) {
	for payload := range client.Subscribe(channel).Channel() {
		var data TokenChangeData
		if err := json.Unmarshal([]byte(payload.Payload), &data); err != nil {
			continue
		}

		ch <- data
	}
}
