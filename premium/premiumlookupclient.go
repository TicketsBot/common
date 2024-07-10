package premium

import (
	"context"
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis/v8"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/ratelimit"
)

type IPremiumLookupClient interface {
	GetCachedTier(ctx context.Context, guildId uint64) (CachedTier, error)
	SetCachedTier(ctx context.Context, guildId uint64, tier CachedTier) error
	DeleteCachedTier(ctx context.Context, guildId uint64) error

	GetTierByGuild(ctx context.Context, guild guild.Guild) (PremiumTier, Source, error)
	GetTierByGuildId(ctx context.Context, guildId uint64, includeVoting bool, botToken string, rateLimiter *ratelimit.Ratelimiter) (PremiumTier, error)
	GetTierByGuildIdWithSource(ctx context.Context, guildId uint64, botToken string, rateLimiter *ratelimit.Ratelimiter) (PremiumTier, Source, error)

	GetTierByUser(ctx context.Context, userId uint64, includeVoting bool) (PremiumTier, error)
	GetTierByUserWithSource(ctx context.Context, userId uint64) (PremiumTier, Source, error)
}

type PremiumLookupClient struct {
	patreonClient *PatreonClient
	redis         *redis.Client
	cache         *cache.PgCache
	database      *database.Database
}

var _ IPremiumLookupClient = (*PremiumLookupClient)(nil)

func NewPremiumLookupClient(patreonClient *PatreonClient, redisClient *redis.Client, cache *cache.PgCache, database *database.Database) *PremiumLookupClient {
	return &PremiumLookupClient{
		patreonClient: patreonClient,
		redis:         redisClient,
		cache:         cache,
		database:      database,
	}
}
