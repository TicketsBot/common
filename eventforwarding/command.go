package eventforwarding

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

type Command struct {
	BotToken     string          `json:"bot_token"`
	BotId        uint64          `json:"bot_id"`
	IsWhitelabel bool            `json:"is_whitelabel"`
	Event        json.RawMessage `json:"data"`
}

const commandKey = "tickets:commands"

func ListenCommands(redis *redis.Client) chan Command {
	ch := make(chan Command)

	go func() {
		for {
			data, err := redis.BLPop(0, commandKey).Result()
			if err != nil || len(data) < 2 {
				fmt.Println(err.Error())
				continue
			}

			// data = [list_name, content]
			var command Command
			if err := json.Unmarshal([]byte(data[1]), &command); err != nil {
				fmt.Println(err.Error())
				continue
			}

			ch <- command
		}
	}()

	return ch
}
