package premium

import (
	"encoding/json"
	"fmt"
	"time"
)

type CachedTier struct {
	Tier   int    `json:"tier"`
	Source Source `json:"source"`
}

const timeout = time.Minute * 5

// Functions can take a user ID or guild ID

func (p *PremiumLookupClient) GetCachedTier(id uint64) (tier CachedTier, err error) {
	key := fmt.Sprintf("premium:%d", id)

	res, err := p.redis.Get(key).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(res), &tier)
	return
}

func (p *PremiumLookupClient) SetCachedTier(id uint64, data CachedTier) (err error) {
	key := fmt.Sprintf("premium:%d", id)

	marshalled, err := json.Marshal(data)
	if err != nil {
		return
	}

	return p.redis.Set(key, string(marshalled), timeout).Err()
}
