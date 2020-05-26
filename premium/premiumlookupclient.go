package premium

import (
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis"
	"github.com/rxdn/gdl/cache"
)

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
