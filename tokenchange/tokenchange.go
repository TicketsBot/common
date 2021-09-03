package tokenchange

import (
	"context"
	"encoding/json"
	"github.com/TicketsBot/common/utils"
	"github.com/go-redis/redis/v8"
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

	return client.Publish(utils.DefaultContext(), channel, string(marshalled)).Err()
}

func ListenTokenChange(client *redis.Client, ch chan TokenChangeData) {
	for payload := range client.Subscribe(context.Background(), channel).Channel() {
		var data TokenChangeData
		if err := json.Unmarshal([]byte(payload.Payload), &data); err != nil {
			continue
		}

		ch <- data
	}
}
