package eventforwarding

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

type Event struct {
	BotToken     string          `json:"bot_token"`
	BotId        uint64          `json:"bot_id"`
	IsWhitelabel bool            `json:"is_whitelabel"`
	ShardId      int             `json:"shard_id"`
	Event        json.RawMessage `json:"event"`
}

const key = "tickets:events"

func ForwardEvent(redis *redis.Client, data Event) error {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return redis.RPush(key, string(marshalled)).Err()
}

func Listen(redis *redis.Client) chan Event {
	ch := make(chan Event)

	go func() {
		for {
			data, err := redis.BLPop(0, key).Result()
			if err != nil || len(data) < 2 {
				fmt.Println(err.Error())
				continue
			}

			// data = [list_name, content]
			var event Event
			if err := json.Unmarshal([]byte(data[1]), &event); err != nil {
				fmt.Println(err.Error())
				continue
			}

			ch <- event
		}
	}()

	return ch
}
