package permission

import (
	"github.com/TicketsBot/database"
	"github.com/go-redis/redis"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/permission"
)

type Retriever interface {
	Db() *database.Database
	Redis() *redis.Client
	IsBotAdmin(uint64) bool
	GetGuild(uint64) (guild.Guild, error)
	GetChannel(uint64) (channel.Channel, error)
	GetGuildMember(guildId, userId uint64) (member.Member, error)
	GetGuildRoles(uint64) ([]guild.Role, error)
}

func GetPermissionLevel(retriever Retriever, member member.Member, guildId uint64) PermissionLevel {
	// Check user ID in cache
	if cached, found := GetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id); found {
		return cached
	}

	// Check if the user is a bot admin user
	if retriever.IsBotAdmin(member.User.Id) {
		return Admin
	}

	// Check if user is guild owner
	guild, err := retriever.GetGuild(guildId)
	if err == nil {
		if member.User.Id == guild.OwnerId {
			go SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Admin)
			return Admin
		}
	}

	// Check user perms for admin
	adminUser, _ := retriever.Db().Permissions.IsAdmin(guildId, member.User.Id)
	if adminUser {
		go SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Admin)
		return Admin
	}

	// Check roles from DB
	adminRoles, _ := retriever.Db().RolePermissions.GetAdminRoles(guildId)
	for _, adminRoleId := range adminRoles {
		if member.HasRole(adminRoleId) {
			go SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Admin)
			return Admin
		}
	}

	// Check if user has Administrator permission
	hasAdminPermission := HasPermissions(retriever, guildId, member.User.Id, permission.Administrator)
	if hasAdminPermission {
		go SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Admin)
		return Admin
	}

	// Check user perms for support
	isSupport, _ := retriever.Db().Permissions.IsSupport(guildId, member.User.Id)
	if isSupport {
		go SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Support)
		return Support
	}

	// Check DB for support roles
	supportRoles, _ := retriever.Db().RolePermissions.GetSupportRoles(guildId)
	for _, supportRoleId := range supportRoles {
		if member.HasRole(supportRoleId) {
			go SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Support)
			return Support
		}
	}

	go SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Everyone)
	return Everyone
}
