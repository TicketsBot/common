package premium

import (
	"time"
)

func (p *PremiumLookupClient) GetTierByUser(userId uint64, includeVoting bool) (tier PremiumTier) {
	// check for cached result
	cached, err := p.getCachedTier(userId)
	if err == nil {
		if includeVoting || !cached.FromVoting {
			if tier = PremiumTier(cached.Tier); tier > None {
				return
			}
		}
	}

	var fromVoting bool

	defer func() {
		// cache result
		go p.setCachedTier(userId, cachedTier{
			Tier:       int(tier),
			FromVoting: fromVoting,
		})
	}()

	// check patreon
	if tier, err = p.patreonClient.getTier(userId); tier > None && err == nil {
		return
	}

	// check for votes
	if includeVoting {
		if tier, err = p.hasVoted(userId); tier > None && err == nil {
			fromVoting = true
			return
		}
	}

	return None
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

