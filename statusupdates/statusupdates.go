package statusupdates

import (
	"github.com/go-redis/redis"
	"strconv"
)

const channel = "tickets:statusupdates"

func Publish(redis *redis.Client, botId uint64) {
	redis.Publish(channel, botId)
}

func Listen(redis *redis.Client, ch chan uint64) {
	for payload := range redis.Subscribe(channel).Channel() {
		if id, err := strconv.ParseUint(payload.Payload, 10, 64); err == nil {
			ch <- id
		}
	}
}
