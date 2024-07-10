package permission

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const redisTimeout = time.Minute * 5

type RedisCache struct {
	client *redis.Client
}

var _ PermissionCache = (*RedisCache)(nil)

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) GetCachedPermissionLevel(ctx context.Context, guildId, userId uint64) (PermissionLevel, error) {
	key := fmt.Sprintf("permissions:%d:%d", guildId, userId)

	res, err := c.client.Get(ctx, key).Result()
	switch {
	case err == nil:
	case errors.Is(err, redis.Nil):
		return Everyone, ErrNotCached
	default:
		return Everyone, err
	}

	parsed, err := strconv.Atoi(res)
	if err != nil {
		return Everyone, err
	}

	return PermissionLevel(parsed), nil
}

func (c *RedisCache) SetCachedPermissionLevel(ctx context.Context, guildId, userId uint64, level PermissionLevel) error {
	key := fmt.Sprintf("permissions:%d:%d", guildId, userId)
	return c.client.Set(ctx, key, level.Int(), redisTimeout).Err()
}

func (c *RedisCache) DeleteCachedPermissionLevel(ctx context.Context, guildId, userId uint64) error {
	key := fmt.Sprintf("permissions:%d:%d", guildId, userId)
	return c.client.Del(ctx, key).Err()
}
