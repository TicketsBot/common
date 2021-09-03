package eventforwarding

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rxdn/gdl/objects/interaction"
)

type Interaction struct {
	BotToken        string                      `json:"bot_token"`
	BotId           uint64                      `json:"bot_id"`
	IsWhitelabel    bool                        `json:"is_whitelabel"`
	InteractionType interaction.InteractionType `json:"interaction_type"`
	Event           json.RawMessage             `json:"data"`
}

const commandKey = "tickets:commands"

func ListenCommands(redis *redis.Client) chan Interaction {
	ch := make(chan Interaction)

	go func() {
		for {
			data, err := redis.BLPop(context.Background(), 0, commandKey).Result()
			if err != nil || len(data) < 2 {
				fmt.Println(err.Error())
				continue
			}

			// data = [list_name, content]
			var command Interaction
			if err := json.Unmarshal([]byte(data[1]), &command); err != nil {
				fmt.Println(err.Error())
				continue
			}

			ch <- command
		}
	}()

	return ch
}
