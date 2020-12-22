package permission

import "errors"

var ErrNotCached = errors.New("member permission level is not cached")

type PermissionCache interface {
	GetCachedPermissionLevel(guildId, userId uint64) (PermissionLevel, error)
	SetCachedPermissionLevel(guildId, userId uint64, level PermissionLevel) error
}
