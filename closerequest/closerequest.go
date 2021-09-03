package closerequest

import (
	"context"
	"encoding/json"
	"github.com/TicketsBot/common/utils"
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis/v8"
)

const channel = "tickets:closerequest:timer"

func PublishMessage(redis *redis.Client, data database.CloseRequest) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return redis.RPush(utils.DefaultContext(), channel, string(json)).Err()
}

func Listen(redis *redis.Client, ch chan database.CloseRequest) {
	for {
		data, err := redis.BLPop(context.Background(), 0, channel).Result()
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
