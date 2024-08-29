package model

import (
	"github.com/google/uuid"
	"time"
)

type Entitlement struct {
	Id        uuid.UUID         `json:"id"`
	GuildId   *uint64           `json:"guild_id"`
	UserId    *uint64           `json:"user_id"`
	SkuId     uuid.UUID         `json:"sku_id"`
	Source    EntitlementSource `json:"source"`
	ExpiresAt *time.Time        `json:"expires_at"`
}

type EntitlementSource string

const (
	EntitlementSourceDiscord EntitlementSource = "discord"
	EntitlementSourcePatreon EntitlementSource = "patreon"
	EntitlementSourceVoting  EntitlementSource = "voting"
	EntitlementSourceKey     EntitlementSource = "key"
)

type EntitlementTier string

const (
	EntitlementTierPremium    EntitlementTier = "premium"
	EntitlementTierWhitelabel EntitlementTier = "whitelabel"
)

type Sku struct {
	Id      uuid.UUID `json:"id"`
	Label   string    `json:"label"`
	SkuType SkuType   `json:"sku_type"`
}

type SkuType string

const (
	SkuTypeSubscription SkuType = "subscription"
	SkuTypeConsumable   SkuType = "consumable"
	SkuTypeDurable      SkuType = "durable"
)

type SubscriptionSku struct {
	SkuId    uuid.UUID       `json:"sku_id"`
	Tier     EntitlementTier `json:"tier"`
	Priority int32           `json:"priority"`
	IsGlobal bool            `json:"is_global"`
}

type WhitelabelSkuData struct {
	SkuId                  uuid.UUID `json:"sku_id"`
	BotPermitted           int       `json:"bot_permitted"`
	ServersPerBotPermitted *int      `json:"servers_per_bot_permitted"`
}

type GuildEntitlementEntry struct {
	Id          uuid.UUID         `json:"id"`
	UserId      uint64            `json:"user_id"`
	Source      EntitlementSource `json:"source"`
	ExpiresAt   *time.Time        `json:"expires_at"`
	SkuId       uuid.UUID         `json:"sku_id"`
	SkuLabel    string            `json:"sku_label"`
	Tier        EntitlementTier   `json:"tier"`
	SkuPriority int32             `json:"sku_priority"`
}
