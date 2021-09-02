package autoclose

import (
	"encoding/json"
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis"
)

const channel = "tickets:closerequest:timer"

func PublishMessage(redis *redis.Client, data []database.CloseRequest) error {
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
