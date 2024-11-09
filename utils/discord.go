package utils

import (
	"time"
)

const DiscordEpoch = 1420070400000

func SnowflakeToTimestamp(snowflake uint64) time.Time {
	epoch := (snowflake >> uint64(22)) + DiscordEpoch
	return time.Unix(int64(epoch)/1000, 0)
}
