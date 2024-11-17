package premium

import (
	"context"
	"errors"
	"github.com/TicketsBot/common/model"
	"github.com/TicketsBot/common/sentry"
	"github.com/go-redis/redis/v8"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/guild"
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
	return sentry.WithSpan3(ctx, "GetTierByGuildIdWithSource", func(span *sentry.Span) (PremiumTier, model.EntitlementSource, error) {
		// check for cached tier by guild ID
		cached, err := p.GetCachedTier(ctx, guildId)
		if err != nil && err != redis.Nil {
			return None, "", err
		} else if err == nil {
			return PremiumTier(cached.Tier), cached.Source, nil
		}

		// retrieve guild object
		g, err := sentry.WithSpan2(ctx, "GetCachedGuild", func(span *sentry.Span) (guild.Guild, error) {
			return p.cache.GetGuild(ctx, guildId)
		})
		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return None, "", err
		}

		if errors.Is(err, cache.ErrNotFound) || g.OwnerId == 0 {
			g, err = sentry.WithSpan2(ctx, "GetGuild", func(span *sentry.Span) (guild.Guild, error) {
				return rest.GetGuild(ctx, botToken, ratelimiter, guildId)
			})
			if err != nil {
				return None, "", err
			}

			go p.cache.StoreGuild(ctx, g)
		}

		return p.GetTierByGuild(ctx, g)
	})
}
