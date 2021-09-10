package permission

import (
	"sync"
	"time"
)

const memoryTimeout = time.Minute

type MemoryCache struct {
	store map[memberId]PermissionLevel
	mu    sync.RWMutex

	cancelRemoval   map[memberId]chan struct{}
	cancelRemovalMu sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store:         make(map[memberId]PermissionLevel),
		cancelRemoval: make(map[memberId]chan struct{}),
	}
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

	c.cancelRemovalMu.Lock()

	c.mu.Lock()
	c.store[member] = level
	c.mu.Unlock()

	cancel := make(chan struct{})

	if existing, ok := c.cancelRemoval[member]; ok {
		existing <- struct{}{}
	}

	c.cancelRemoval[member] = cancel
	c.cancelRemovalMu.Unlock()

	timer := time.NewTimer(memoryTimeout)
	go func() {
		select {
		case <-timer.C:
			c.cancelRemovalMu.Lock()
			delete(c.cancelRemoval, member)

			c.mu.Lock()
			delete(c.store, member)
			c.mu.Unlock()

			c.cancelRemovalMu.Unlock()
		case <-cancel:
		}
	}()

	return nil
}

func (c *MemoryCache) DeleteCachedPermissionLevel(guildId, userId uint64) error {
	member := memberId{guildId, userId}

	c.mu.Lock()
	delete(c.store, member)
	c.mu.Unlock()

	if existing, ok := c.cancelRemoval[member]; ok {
		existing <- struct{}{}
	}

	c.cancelRemovalMu.Lock()
	delete(c.cancelRemoval, member)
	c.cancelRemovalMu.Unlock()

	return nil
}
