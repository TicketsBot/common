package premium

import (
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/ratelimit"
)

type MockLookupClient struct {
	Tier   PremiumTier
	Source Source
}

func NewMockLookupClient(tier PremiumTier, src Source) MockLookupClient {
	return MockLookupClient{
		Tier:   tier,
		Source: src,
	}
}

func (c *MockLookupClient) GetCachedTier(uint64) (CachedTier, error) {
	return CachedTier{
		Tier:   int8(c.Tier),
		Source: c.Source,
	}, nil
}

func (c *MockLookupClient) SetCachedTier(uint64, CachedTier) error {
	return nil
}

func (c *MockLookupClient) GetTierByGuild(guild.Guild) (PremiumTier, Source, error) {
	return c.Tier, c.Source, nil
}

func (c *MockLookupClient) GetTierByGuildId(_ uint64, includeVoting bool, _ string, _ *ratelimit.Ratelimiter) (PremiumTier, error) {
	if !includeVoting && c.Source == SourceVoting {
		return None, nil
	}

	return c.Tier, nil
}

func (c *MockLookupClient) GetTierByGuildIdWithSource(uint64, string, *ratelimit.Ratelimiter) (PremiumTier, Source, error) {
	return c.Tier, c.Source, nil
}

func (c *MockLookupClient) GetTierByUser(_ uint64, includeVoting bool) (PremiumTier, error) {
	if !includeVoting && c.Source == SourceVoting {
		return None, nil
	}

	return c.Tier, nil
}

func (c *MockLookupClient) GetTierByUserWithSource(uint64) (PremiumTier, Source, error) {
	return c.Tier, c.Source, nil
}
