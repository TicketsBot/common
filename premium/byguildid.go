package premium

import (
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/ratelimit"
)

func (p *PremiumLookupClient) GetTierByGuildId(guildId uint64, includeVoting bool, botToken string, ratelimiter *ratelimit.Ratelimiter) (tier PremiumTier) {
	// retrieve guild object
	guild, found := p.cache.GetGuild(guildId, false)
	if !found {
		var err error
		guild, err = rest.GetGuild(botToken, ratelimiter, guildId)
		if err != nil {
			return None
		}

		go p.cache.StoreGuild(guild)
	}

	return p.GetTierByGuild(guild, includeVoting)
}
