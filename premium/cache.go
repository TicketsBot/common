package premium

import (
	"encoding/json"
	"fmt"
	"time"
)

type cachedTier struct {
	Tier       int    `json:"tier"`
	FromVoting bool   `json:"from_voting"`
}

const timeout = time.Minute * 5

func (p *PremiumLookupClient) getCachedTier(userId uint64) (tier cachedTier, err error) {
	key := fmt.Sprintf("premium:%d", userId)

	res, err := p.redis.Get(key).Result(); if err != nil {
		return
	}

	err = json.Unmarshal([]byte(res), &tier)
	return
}

func (p *PremiumLookupClient) setCachedTier(userId uint64, data cachedTier) (err error) {
	key := fmt.Sprintf("premium:%d", userId)

	marshalled, err := json.Marshal(data); if err != nil {
		return
	}

	return p.redis.Set(key, string(marshalled), timeout).Err()
}
