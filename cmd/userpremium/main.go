package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rxdn/gdl/cache"
	"os"
)

var userId = flag.Uint64("userid", 0, "user id to check")

func main() {
	flag.Parse()

	client := createClient()
	tier, src, err := client.GetTierByUserWithSource(context.Background(), *userId)
	must(err)

	fmt.Printf("%s via %s\n", tier, src)
}

func createClient() *premium.PremiumLookupClient {
	patreonClient := premium.NewPatreonClient(os.Getenv("PATREON_URL"), os.Getenv("PATREON_KEY"))

	redisClient := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	cachePool, err := pgxpool.Connect(context.Background(), os.Getenv("CACHE_URI"))
	must(err)

	cacheClient := cache.NewPgCache(cachePool, cache.CacheOptions{
		Guilds:      true,
		Users:       true,
		Members:     true,
		Channels:    true,
		Roles:       true,
		Emojis:      false,
		VoiceStates: false,
	})

	dbPool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URI"))
	must(err)

	dbClient := database.NewDatabase(dbPool)

	return premium.NewPremiumLookupClient(patreonClient, redisClient, &cacheClient, dbClient)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
