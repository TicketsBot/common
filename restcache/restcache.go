package restcache

import "github.com/rxdn/gdl/objects/guild"

type RestCache interface {
	GetGuildRoles(guildId uint64) ([]guild.Role, error)
}
