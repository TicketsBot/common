package premium

import (
	"context"
	"github.com/TicketsBot/common/model"
	"github.com/go-redis/redis/v8"
)

func (p *PremiumLookupClient) GetTierByUser(ctx context.Context, userId uint64, includeVoting bool) (PremiumTier, error) {
	tier, source, err := p.GetTierByUserWithSource(ctx, userId)
	if err != nil {
		return None, err
	}

	if source == model.EntitlementSourceVoting && !includeVoting {
		return None, nil
	}

	return tier, nil
}

func (p *PremiumLookupClient) GetTierByUserWithSource(ctx context.Context, userId uint64) (_tier PremiumTier, _src model.EntitlementSource, _err error) {
	_tier = None
	_src = ""

	// check for cached result
	cached, err := p.GetCachedTier(ctx, userId)
	if err != nil && err != redis.Nil {
		return None, "", err
	} else if err == nil {
		return PremiumTier(cached.Tier), cached.Source, nil
	}

	defer func() {
		// cache result
		if _err == nil {
			go p.SetCachedTier(ctx, userId, CachedTier{
				Tier:   int8(_tier),
				Source: _src,
			})
		}
	}()

	// check entitlements db
	subscriptions, err := p.database.Entitlements.ListUserSubscriptions(ctx, userId, GracePeriod)
	if err != nil {
		return None, "", err
	}

	if maxSubscription := findMaxTier(subscriptions); maxSubscription != nil {
		return TierFromEntitlement(maxSubscription.Tier), maxSubscription.Source, nil
	}

	// check for votes
	votingTier, err := p.hasVoted(ctx, userId)
	if err != nil {
		return None, "", err
	} else if votingTier > None {
		return votingTier, model.EntitlementSourceVoting, nil
	}

	return None, "", nil
}

func (p *PremiumLookupClient) getTierByUsers(ctx context.Context, userIds []uint64) (tier PremiumTier, src model.EntitlementSource, _err error) {
	tier = None
	src = ""

	// check for votes
	// we can skip here if already premium
	if tier == None {
		votingTier, err := p.hasVoted(ctx, userIds...)
		if err != nil {
			return None, "", err
		} else if votingTier > tier {
			return votingTier, model.EntitlementSourceVoting, nil
		}
	}

	return
}

func (p *PremiumLookupClient) hasVoted(ctx context.Context, userIds ...uint64) (PremiumTier, error) {
	isPremium, err := p.database.Votes.Any(ctx, userIds...)
	if err != nil {
		return None, err
	}

	if isPremium {
		return Premium, nil
	} else {
		return None, err
	}
}
