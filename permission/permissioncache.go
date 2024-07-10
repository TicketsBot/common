package permission

import (
	"context"
	"errors"
)

var ErrNotCached = errors.New("member permission level is not cached")

type PermissionCache interface {
	GetCachedPermissionLevel(ctx context.Context, guildId, userId uint64) (PermissionLevel, error)
	SetCachedPermissionLevel(ctx context.Context, guildId, userId uint64, level PermissionLevel) error
	DeleteCachedPermissionLevel(ctx context.Context, guildId, userId uint64) error
}
