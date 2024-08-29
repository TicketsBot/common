package premium

import (
	"context"
	"errors"
	"github.com/TicketsBot/common/model"
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

	if source == model.EntitlementSourceVoting && !includeVoting {
		return None, nil
	}

	return tier, nil
}

func (p *PremiumLookupClient) GetTierByGuildIdWithSource(ctx context.Context, guildId uint64, botToken string, ratelimiter *ratelimit.Ratelimiter) (PremiumTier, model.EntitlementSource, error) {
	// check for cached tier by guild ID
	cached, err := p.GetCachedTier(ctx, guildId)
	if err != nil && err != redis.Nil {
		return None, "", err
	} else if err == nil {
		return PremiumTier(cached.Tier), cached.Source, nil
	}

	// retrieve guild object
	guild, err := p.cache.GetGuild(ctx, guildId)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		return None, "", err
	}

	if errors.Is(err, cache.ErrNotFound) || guild.OwnerId == 0 {
		guild, err = rest.GetGuild(ctx, botToken, ratelimiter, guildId)
		if err != nil {
			return None, "", err
		}

		go p.cache.StoreGuild(ctx, guild)
	}

	return p.GetTierByGuild(ctx, guild)
}
