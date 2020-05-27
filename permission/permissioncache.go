package permission

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

const timeout = time.Minute * 5

func GetCachedPermissionLevel(redis *redis.Client, guildId, userId uint64) (PermissionLevel, bool) {
	key := fmt.Sprintf("permissions:%d:%d", guildId, userId)
	res, err := redis.Get(key).Result(); if err != nil {
		return Everyone, false
	}

	parsed, err := strconv.Atoi(res); if err != nil {
		return Everyone, false
	}

	return PermissionLevel(parsed), true
}

func SetCachedPermissionLevel(redis *redis.Client, guildId, userId uint64, level PermissionLevel) error {
	key := fmt.Sprintf("permissions:%d:%d", guildId, userId)
	return redis.Set(key, level.Int(), timeout).Err()
}
