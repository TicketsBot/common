package premium

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/common/model"
	"time"
)

type CachedTier struct {
	Tier   int8                    `json:"tier"`
	Source model.EntitlementSource `json:"source"`
}

const timeout = time.Minute * 5

// Functions can take a user ID or guild ID

func (p *PremiumLookupClient) GetCachedTier(ctx context.Context, id uint64) (tier CachedTier, err error) {
	key := fmt.Sprintf("premium:%d", id)

	res, err := p.redis.Get(ctx, key).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(res), &tier)
	return
}

func (p *PremiumLookupClient) SetCachedTier(ctx context.Context, id uint64, data CachedTier) (err error) {
	key := fmt.Sprintf("premium:%d", id)

	marshalled, err := json.Marshal(data)
	if err != nil {
		return
	}

	return p.redis.Set(ctx, key, string(marshalled), timeout).Err()
}

func (p *PremiumLookupClient) DeleteCachedTier(ctx context.Context, id uint64) (err error) {
	key := fmt.Sprintf("premium:%d", id)
	return p.redis.Del(ctx, key).Err()
}
