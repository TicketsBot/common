package restcache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/ratelimit"
	"time"
)

const (
	LockTimeout      = 5 * time.Second
	LockRetryTimeout = 10 * time.Second
	RedisTimeout     = 3 * time.Second
	RetryDelay       = 100 * time.Millisecond
	CacheExpiry      = time.Minute
)

var ErrFailedToAcquireLock = fmt.Errorf("failed to acquire lock")

type RedisRestCache struct {
	client      *redis.Client
	botToken    string
	ratelimiter *ratelimit.Ratelimiter
}

func NewRedisRestCache(client *redis.Client, botToken string, ratelimiter *ratelimit.Ratelimiter) *RedisRestCache {
	return &RedisRestCache{
		client,
		botToken,
		ratelimiter,
	}
}

func (c *RedisRestCache) GetGuildRoles(guildId uint64) ([]guild.Role, error) {
	context, cancel := context.WithTimeout(context.Background(), LockTimeout)
	defer cancel()

	var roles []guild.Role
	err := getOrFetch(c, context, fmt.Sprintf("roles:%d", guildId), &roles, func() ([]guild.Role, error) {
		return rest.GetGuildRoles(c.botToken, c.ratelimiter, guildId)
	})

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func getOrFetch[T any](c *RedisRestCache, context context.Context, key string, v *T, populator func() (T, error)) (finalErr error) {
	value, err := c.client.Get(context, key).Result()
	if err != nil {
		if err == redis.Nil { // key doesn't exist
			lockKey := fmt.Sprintf("%s:lock", key)

			hasLock := false
			i := 0
			for !hasLock {
				hasLock, err = c.client.SetNX(context, lockKey, "1", LockTimeout).Result()
				if err != nil {
					return err
				}

				i++
				if i > 15 {
					return ErrFailedToAcquireLock
				}

				if !hasLock {
					time.Sleep(RetryDelay)
				}
			}

			defer func() {
				c.client.Del(context, lockKey) // Context should prevent us form unlocking a lock that has expired?
			}()

			value, err := populator()
			if err != nil {
				return err
			}

			*v = value

			encoded, err := json.Marshal(value)
			if err != nil {
				return err
			}

			if err := c.client.Set(context, key, encoded, CacheExpiry).Err(); err != nil {
				return err
			}

			return nil
		} else {
			return err
		}
	}

	return json.Unmarshal([]byte(value), &v)
}
