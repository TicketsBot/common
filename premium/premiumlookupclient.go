package premium

import (
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis/v8"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/ratelimit"
)

type IPremiumLookupClient interface {
	GetCachedTier(uint64) (CachedTier, error)
	SetCachedTier(uint64, CachedTier) error
	DeleteCachedTier(uint64) error

	GetTierByGuild(guild.Guild) (PremiumTier, Source, error)
	GetTierByGuildId(uint64, bool, string, *ratelimit.Ratelimiter) (PremiumTier, error)
	GetTierByGuildIdWithSource(uint64, string, *ratelimit.Ratelimiter) (PremiumTier, Source, error)

	GetTierByUser(uint64, bool) (PremiumTier, error)
	GetTierByUserWithSource(uint64) (PremiumTier, Source, error)
}

type PremiumLookupClient struct {
	patreonClient *PatreonClient
	redis         *redis.Client
	cache         *cache.PgCache
	database      *database.Database
}

func NewPremiumLookupClient(patreonClient *PatreonClient, redisClient *redis.Client, cache *cache.PgCache, database *database.Database) *PremiumLookupClient {
	return &PremiumLookupClient{
		patreonClient: patreonClient,
		redis:         redisClient,
		cache:         cache,
		database:      database,
	}
}
