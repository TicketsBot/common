package premium

import (
	"github.com/TicketsBot/common/sentry"
	"github.com/rxdn/gdl/objects/guild"
)

func (p *PremiumLookupClient) GetTierByGuild(guild guild.Guild) (_tier PremiumTier, _src Source, _err error) {
	_tier = None
	_src = -1

	defer func() {
		// cache result
		if _err == nil {
			go func() {
				err := p.SetCachedTier(guild.Id, CachedTier{
					Tier:   int8(_tier),
					Source: _src,
				})

				if err != nil {
					sentry.Error(err)
				}
			}()
		}
	}()

	admins, err := p.database.Permissions.GetAdmins(guild.Id)
	if err != nil {
		return None, -1, err
	}

	admins = append(admins, guild.OwnerId)

	// check patreon + votes
	// key lookup cannot be whitelabel, therefore we don't need to do key lookup if patreon is regular premium or higher
	adminsTier, src, err := p.getTierByUsers(admins)
	if err != nil {
		return None, -1, err
	} else if adminsTier > None {
		return adminsTier, src, nil
	}

	keyTier, err := p.hasKey(guild.Id)
	if err != nil {
		return None, -1, err
	} else if keyTier > None {
		return keyTier, SourcePremiumKey, nil
	}

	return None, -1, nil
}

func (p *PremiumLookupClient) hasKey(guildId uint64) (PremiumTier, error) {
	isPremium, err := p.database.PremiumGuilds.IsPremium(guildId)
	if err != nil {
		return None, err
	}

	if isPremium {
		return Premium, nil
	} else {
		return None, nil
	}
}
