package premium

import "github.com/rxdn/gdl/objects/guild"

func (p *PremiumLookupClient) GetTierByGuild(guild guild.Guild, includeVoting bool) (tier PremiumTier) {
	// check patreon + votes
	// key lookup cannot be whitelabel, therefore we don't need to do key lookup if patreon is regular premium or higher
	if tier = p.GetTierByUser(guild.OwnerId, includeVoting); tier > None {
		return
	}

	var err error
	if tier, err = p.hasKey(guild.Id); err == nil && tier > None {
		return
	}

	return None
}

func (p *PremiumLookupClient) hasKey(guildId uint64) (PremiumTier, error) {
	isPremium, err := p.database.PremiumGuilds.IsPremium(guildId); if err != nil {
		return None, err
	}

	if isPremium {
		return Premium, nil
	} else {
		return None, nil
	}
}
