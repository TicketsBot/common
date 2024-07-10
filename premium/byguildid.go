package premium

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/ratelimit"
)

func (p *PremiumLookupClient) GetTierByGuildId(ctx context.Context, guildId uint64, includeVoting bool, botToken string, ratelimiter *ratelimit.Ratelimiter) (PremiumTier, error) {
	tier, source, err := p.GetTierByGuildIdWithSource(ctx, guildId, botToken, ratelimiter)
	if err != nil {
		return None, err
	}

	if source == SourceVoting && !includeVoting {
		return None, nil
	}

	return tier, nil
}

func (p *PremiumLookupClient) GetTierByGuildIdWithSource(ctx context.Context, guildId uint64, botToken string, ratelimiter *ratelimit.Ratelimiter) (PremiumTier, Source, error) {
	// check for cached tier by guild ID
	cached, err := p.GetCachedTier(ctx, guildId)
	if err != nil && err != redis.Nil {
		return None, -1, err
	} else if err == nil {
		return PremiumTier(cached.Tier), cached.Source, nil
	}

	// retrieve guild object
	guild, err := p.cache.GetGuild(ctx, guildId)
	if err == nil {
		var err error
		guild, err = rest.GetGuild(ctx, botToken, ratelimiter, guildId)
		if err != nil {
			return None, -1, err
		}

		go p.cache.StoreGuild(ctx, guild)
	} else if !errors.Is(err, cache.ErrNotFound) {
		return None, -1, err
	}

	return p.GetTierByGuild(ctx, guild)
}
