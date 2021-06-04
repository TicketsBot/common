package premium

import "github.com/rxdn/gdl/objects/guild"

func (p *PremiumLookupClient) GetTierByGuild(guild guild.Guild, includeVoting bool) (tier PremiumTier) {
	var fromVoting bool

	defer func() {
		go p.SetCachedTier(guild.Id, CachedTier{
			Tier:       int(tier),
			FromVoting: fromVoting,
		})
	}()

	admins, err := p.database.Permissions.GetAdmins(guild.Id)
	if err != nil { // TODO: LOG
		return None
	}

	admins = append(admins, guild.OwnerId)

	// check patreon + votes
	// key lookup cannot be whitelabel, therefore we don't need to do key lookup if patreon is regular premium or higher
	if tier, fromVoting = p.getTierByUsers(admins, includeVoting); tier > None {
		return
	}

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
