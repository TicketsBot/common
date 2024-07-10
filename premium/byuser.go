package premium

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func (p *PremiumLookupClient) GetTierByUser(ctx context.Context, userId uint64, includeVoting bool) (PremiumTier, error) {
	tier, source, err := p.GetTierByUserWithSource(ctx, userId)
	if err != nil {
		return None, err
	}

	if source == SourceVoting && !includeVoting {
		return None, nil
	}

	return tier, nil
}

func (p *PremiumLookupClient) GetTierByUserWithSource(ctx context.Context, userId uint64) (_tier PremiumTier, _src Source, _err error) {
	_tier = None
	_src = -1

	// check for cached result
	cached, err := p.GetCachedTier(ctx, userId)
	if err != nil && err != redis.Nil {
		return None, -1, err
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

	// check patreon
	patreonTier, err := p.patreonClient.GetTier(ctx, userId)
	if err != nil {
		return None, -1, err
	} else if patreonTier > None {
		return patreonTier, SourcePatreon, nil
	}

	// check whitelabel keys
	isWhitelabel, err := p.hasWhitelabelKey(ctx, userId)
	if err != nil {
		return None, -1, err
	} else if isWhitelabel {
		return Whitelabel, SourceWhitelabelKey, nil
	}

	// check for votes
	votingTier, err := p.hasVoted(ctx, userId)
	if err != nil {
		return None, -1, err
	} else if votingTier > None {
		return votingTier, SourceVoting, nil
	}

	return None, -1, nil
}

func (p *PremiumLookupClient) getTierByUsers(ctx context.Context, userIds []uint64) (tier PremiumTier, src Source, _err error) {
	tier = None
	src = -1

	// check patreon
	patreonTier, err := p.patreonClient.GetTier(ctx, userIds...)
	if err != nil {
		return None, -1, err
	} else if patreonTier > tier {
		tier = patreonTier
		src = SourcePatreon
	}

	if tier == Whitelabel {
		return
	}

	// check whitelabel keys
	isWhitelabel, err := p.hasWhitelabelKey(ctx, userIds...)
	if err != nil {
		return None, -1, err
	} else if isWhitelabel {
		return Whitelabel, SourceWhitelabelKey, nil
	}
	// check for votes
	// we can skip here if already premium
	if tier == None {
		votingTier, err := p.hasVoted(ctx, userIds...)
		if err != nil {
			return None, -1, err
		} else if votingTier > tier {
			return votingTier, SourceVoting, nil
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

func (p *PremiumLookupClient) hasWhitelabelKey(ctx context.Context, userIds ...uint64) (bool, error) {
	return p.database.WhitelabelUsers.AnyPremium(ctx, userIds)
}
