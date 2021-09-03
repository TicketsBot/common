package permission

import (
	"fmt"
	"github.com/TicketsBot/common/utils"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const redisTimeout = time.Minute * 5

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) GetCachedPermissionLevel(guildId, userId uint64) (PermissionLevel, error) {
	key := fmt.Sprintf("permissions:%d:%d", guildId, userId)

	res, err := c.client.Get(utils.DefaultContext(), key).Result()
	switch err {
	case nil:
	case redis.Nil:
		return Everyone, ErrNotCached
	default:
		return Everyone, err
	}

	parsed, err := strconv.Atoi(res); if err != nil {
		return Everyone, err
	}

	return PermissionLevel(parsed), nil
}

func (c *RedisCache) SetCachedPermissionLevel(guildId, userId uint64, level PermissionLevel) error {
	key := fmt.Sprintf("permissions:%d:%d", guildId, userId)
	return c.client.Set(utils.DefaultContext(), key, level.Int(), redisTimeout).Err()
}
