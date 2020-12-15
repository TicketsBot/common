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

func GetPermissionLevel(retriever Retriever, member member.Member, guildId uint64) (PermissionLevel, error) {
	// Check user ID in cache
	if cached, err := GetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id); err == nil {
		return cached, nil
	} else if err != redis.Nil {
		return Everyone, err
	}

	// Check if the user is a bot admin user
	if retriever.IsBotAdmin(member.User.Id) {
		return Admin, nil
	}

	// Check if user is guild owner
	guild, err := retriever.GetGuild(guildId)
	if err == nil {
		if member.User.Id == guild.OwnerId {
			err := SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Admin)
			return Admin, err
		}
	}

	// Check user perms for admin
	adminUser, _ := retriever.Db().Permissions.IsAdmin(guildId, member.User.Id)
	if adminUser {
		err := SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Admin)
		return Admin, err
	}

	// Check roles from DB
	adminRoles, _ := retriever.Db().RolePermissions.GetAdminRoles(guildId)
	for _, adminRoleId := range adminRoles {
		if member.HasRole(adminRoleId) {
			err := SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Admin)
			return Admin, err
		}
	}

	// Check if user has Administrator permission
	hasAdminPermission := HasPermissions(retriever, guildId, member.User.Id, permission.Administrator)
	if hasAdminPermission {
		err := SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Admin)
		return Admin, err
	}

	// Check user perms for support
	isSupport, _ := retriever.Db().Permissions.IsSupport(guildId, member.User.Id)
	if isSupport {
		err := SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Support)
		return Support, err
	}

	// Check DB for support roles
	supportRoles, _ := retriever.Db().RolePermissions.GetSupportRoles(guildId)
	for _, supportRoleId := range supportRoles {
		if member.HasRole(supportRoleId) {
			err := SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Support)
			return Support, err
		}
	}

	err := SetCachedPermissionLevel(retriever.Redis(), guildId, member.User.Id, Everyone)
	return Everyone, err
}
