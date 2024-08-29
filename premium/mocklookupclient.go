package premium

import (
	"context"
	"github.com/TicketsBot/common/model"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/ratelimit"
)

type MockLookupClient struct {
	Tier   PremiumTier
	Source model.EntitlementSource
}

var _ IPremiumLookupClient = (*MockLookupClient)(nil)

func NewMockLookupClient(tier PremiumTier, src model.EntitlementSource) MockLookupClient {
	return MockLookupClient{
		Tier:   tier,
		Source: src,
	}
}

func (c *MockLookupClient) GetCachedTier(context.Context, uint64) (CachedTier, error) {
	return CachedTier{
		Tier:   int8(c.Tier),
		Source: c.Source,
	}, nil
}

func (c *MockLookupClient) SetCachedTier(context.Context, uint64, CachedTier) error {
	return nil
}

func (c *MockLookupClient) DeleteCachedTier(context.Context, uint64) error {
	return nil
}

func (c *MockLookupClient) GetTierByGuild(context.Context, guild.Guild) (PremiumTier, model.EntitlementSource, error) {
	return c.Tier, c.Source, nil
}

func (c *MockLookupClient) GetTierByGuildId(_ context.Context, _ uint64, includeVoting bool, _ string, _ *ratelimit.Ratelimiter) (PremiumTier, error) {
	if !includeVoting && c.Source == model.EntitlementSourceVoting {
		return None, nil
	}

	return c.Tier, nil
}

func (c *MockLookupClient) GetTierByGuildIdWithSource(context.Context, uint64, string, *ratelimit.Ratelimiter) (PremiumTier, model.EntitlementSource, error) {
	return c.Tier, c.Source, nil
}

func (c *MockLookupClient) GetTierByUser(_ context.Context, _ uint64, includeVoting bool) (PremiumTier, error) {
	if !includeVoting && c.Source == model.EntitlementSourceVoting {
		return None, nil
	}

	return c.Tier, nil
}

func (c *MockLookupClient) GetTierByUserWithSource(context.Context, uint64) (PremiumTier, model.EntitlementSource, error) {
	return c.Tier, c.Source, nil
}
