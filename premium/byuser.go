package premium

import (
	"time"
)

func (p *PremiumLookupClient) GetTierByUser(userId uint64, includeVoting bool) (tier PremiumTier) {
	tier, _ = p.getTierByUser(userId, includeVoting)
	return
}

func (p *PremiumLookupClient) getTierByUser(userId uint64, includeVoting bool) (tier PremiumTier, fromVoting bool) {
	// check for cached result
	cached, err := p.GetCachedTier(userId)
	if err == nil {
		if includeVoting || !cached.FromVoting {
			if tier = PremiumTier(cached.Tier); tier > None {
				fromVoting = cached.FromVoting
				return
			}
		}
	}

	defer func() {
		// cache result
		go p.SetCachedTier(userId, CachedTier{
			Tier:       int(tier),
			FromVoting: fromVoting,
		})
	}()

	// check patreon
	if tier, err = p.patreonClient.GetTier(userId); tier > None && err == nil {
		return
	}

	// check for votes
	if includeVoting {
		if tier, err = p.hasVoted(userId); tier > None && err == nil {
			fromVoting = true
			return
		}
	}

	return None, false
}

func (p *PremiumLookupClient) hasVoted(userId uint64) (PremiumTier, error) {
	voteTime, err := p.database.Votes.Get(userId); if err != nil {
		return None, err
	}

	if voteTime.After(time.Now().Add(time.Hour * -24)) {
		return Premium, nil
	} else {
		return None, nil
	}
}

