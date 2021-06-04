package premium

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

	// check whitelabel keys
	if isWhitelabel, err := p.hasWhitelabelKey(userId); err == nil && isWhitelabel {
		return Whitelabel, false
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

func (p *PremiumLookupClient) getTierByUsers(userIds []uint64, includeVoting bool) (tier PremiumTier, fromVoting bool) {
	// check patreon
	patreonTier, err := p.patreonClient.GetTier(userIds...)
	if err == nil && patreonTier > tier {
		tier = patreonTier
	}

	if tier == Whitelabel {
		return
	}

	// check whitelabel keys
	if isWhitelabel, err := p.hasWhitelabelKey(userIds...); err == nil && isWhitelabel {
		return Whitelabel, false
	}

	// check for votes
	// we can skip here if
	if includeVoting && tier == None {
		votingTier, err := p.hasVoted(userIds...)
		if err == nil && votingTier > tier {
			tier = votingTier
			fromVoting = true
			return
		}
	}

	return
}

func (p *PremiumLookupClient) hasVoted(userIds ...uint64) (PremiumTier, error) {
	isPremium, err := p.database.Votes.Any(userIds...)
	if err != nil {
		return None, err
	}

	if isPremium {
		return Premium, nil
	} else {
		return None, err
	}
}

func (p *PremiumLookupClient) hasWhitelabelKey(userIds ...uint64) (bool, error) {
	return p.database.WhitelabelUsers.AnyPremium(userIds)
}
