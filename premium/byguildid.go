package premium

import (
	"github.com/go-redis/redis/v8"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/ratelimit"
)

func (p *PremiumLookupClient) GetTierByGuildId(guildId uint64, includeVoting bool, botToken string, ratelimiter *ratelimit.Ratelimiter) (PremiumTier, error) {
	tier, source, err := p.GetTierByGuildIdWithSource(guildId, botToken, ratelimiter)
	if err != nil {
		return None, err
	}

	if source == SourceVoting && !includeVoting {
		return None, nil
	}

	return tier, nil
}

func (p *PremiumLookupClient) GetTierByGuildIdWithSource(guildId uint64, botToken string, ratelimiter *ratelimit.Ratelimiter) (PremiumTier, Source, error) {
	// check for cached tier by guild ID
	cached, err := p.GetCachedTier(guildId)
	if err != nil && err != redis.Nil {
		return None, -1, err
	} else if err == nil {
		return PremiumTier(cached.Tier), cached.Source, nil
	}

	// retrieve guild object
	guild, found := p.cache.GetGuild(guildId, false)
	if !found {
		var err error
		guild, err = rest.GetGuild(botToken, ratelimiter, guildId)
		if err != nil {
			return None, -1, err
		}

		go p.cache.StoreGuild(guild)
	}

	return p.GetTierByGuild(guild)
}
