package permission

import (
	"sync"
	"time"
)

const memoryTimeout = time.Minute

type MemoryCache struct {
	store map[memberId]PermissionLevel
	mu    sync.RWMutex
}

type memberId struct {
	GuildId, UserId uint64
}

func (c *MemoryCache) GetCachedPermissionLevel(guildId, userId uint64) (PermissionLevel, error) {
	member := memberId{guildId, userId}

	c.mu.RLock()
	defer c.mu.RUnlock()

	level, ok := c.store[member]
	if !ok {
		return Everyone, ErrNotCached
	}

	return level, nil
}

func (c *MemoryCache) SetCachedPermissionLevel(guildId, userId uint64, level PermissionLevel) error {
	member := memberId{guildId, userId}

	c.mu.RLock()
	defer c.mu.RUnlock()

	c.store[member] = level

	return nil
}
