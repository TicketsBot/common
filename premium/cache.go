package premium

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/common/model"
	"github.com/TicketsBot/common/sentry"
	"time"
)

type CachedTier struct {
	Tier   int8                    `json:"tier"`
	Source model.EntitlementSource `json:"source"`
}

const timeout = time.Minute * 5

// Functions can take a user ID or guild ID

func (p *PremiumLookupClient) GetCachedTier(ctx context.Context, id uint64) (CachedTier, error) {
	return sentry.WithSpan2(ctx, "GetCachedTier", func(span *sentry.Span) (CachedTier, error) {
		key := fmt.Sprintf("premium:%d", id)

		res, err := p.redis.Get(ctx, key).Result()
		if err != nil {
			return CachedTier{}, err
		}

		var tier CachedTier
		if err := json.Unmarshal([]byte(res), &tier); err != nil {
			return CachedTier{}, err
		}

		return tier, nil
	})
}

func (p *PremiumLookupClient) SetCachedTier(ctx context.Context, id uint64, data CachedTier) error {
	return sentry.WithSpan1(ctx, "SetCachedTier", func(span *sentry.Span) error {
		key := fmt.Sprintf("premium:%d", id)

		marshalled, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return p.redis.Set(ctx, key, string(marshalled), timeout).Err()
	})
}

func (p *PremiumLookupClient) DeleteCachedTier(ctx context.Context, id uint64) error {
	return sentry.WithSpan1(ctx, "DeleteCachedTier", func(span *sentry.Span) error {
		key := fmt.Sprintf("premium:%d", id)
		return p.redis.Del(ctx, key).Err()
	})
}
