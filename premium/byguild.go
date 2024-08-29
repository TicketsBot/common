package premium

import (
	"context"
	"github.com/TicketsBot/common/model"
	"github.com/TicketsBot/common/sentry"
	"github.com/rxdn/gdl/objects/guild"
	"time"
)

const GracePeriod = time.Hour * 24 // TODO: Reduce this to zero?

func (p *PremiumLookupClient) GetTierByGuild(ctx context.Context, guild guild.Guild) (_tier PremiumTier, _src model.EntitlementSource, _err error) {
	_tier = None
	_src = ""

	defer func() {
		// cache result
		if _err == nil {
			go func() {
				err := p.SetCachedTier(ctx, guild.Id, CachedTier{
					Tier:   int8(_tier),
					Source: _src,
				})

				if err != nil {
					sentry.Error(err)
				}
			}()
		}
	}()

	admins, err := p.database.Permissions.GetAdmins(ctx, guild.Id)
	if err != nil {
		return None, "", err
	}

	admins = append(admins, guild.OwnerId)

	// check entitlements db
	subscriptions, err := p.database.Entitlements.ListGuildSubscriptions(ctx, guild.Id, guild.OwnerId, GracePeriod)
	if err != nil {
		return None, "", err
	}

	if maxSubscription := findMaxTier(subscriptions); maxSubscription != nil {
		return TierFromEntitlement(maxSubscription.Tier), maxSubscription.Source, nil
	}

	// check votes + whitelabel keys
	// key lookup cannot be whitelabel, therefore we don't need to do key lookup if patreon is regular premium or higher
	adminsTier, src, err := p.getTierByUsers(ctx, admins)
	if err != nil {
		return None, "", err
	} else if adminsTier > None {
		return adminsTier, src, nil
	}

	keyTier, err := p.hasKey(ctx, guild.Id)
	if err != nil {
		return None, "", err
	} else if keyTier > None {
		return keyTier, model.EntitlementSourceKey, nil
	}

	return None, "", nil
}

func (p *PremiumLookupClient) hasKey(ctx context.Context, guildId uint64) (PremiumTier, error) {
	isPremium, err := p.database.PremiumGuilds.IsPremium(ctx, guildId)
	if err != nil {
		return None, err
	}

	if isPremium {
		return Premium, nil
	} else {
		return None, nil
	}
}

func findMaxTier(subscriptions []model.GuildEntitlementEntry) *model.GuildEntitlementEntry {
	if len(subscriptions) == 0 {
		return nil
	}

	maxTier := subscriptions[0]
	for _, entry := range subscriptions[1:] {
		if entry.SkuPriority > maxTier.SkuPriority {
			maxTier = entry
		}
	}

	return &maxTier
}
