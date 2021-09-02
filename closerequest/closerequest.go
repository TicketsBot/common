package closerequest

import (
	"encoding/json"
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis"
)

const channel = "tickets:closerequest:timer"

func PublishMessage(redis *redis.Client, data database.CloseRequest) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return redis.RPush(channel, string(json)).Err()
}

func Listen(redis *redis.Client, ch chan database.CloseRequest) {
	for {
		data, err := redis.BLPop(0, channel).Result()
		if err != nil || len(data) < 2 {
			continue
		}

		// data = [list_name, content]
		var ticket database.CloseRequest
		if err := json.Unmarshal([]byte(data[1]), &ticket); err != nil {
			continue
		}

		ch <- ticket
	}
}
