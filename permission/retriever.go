package permission

import (
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/permission"
)

type Retriever interface {
	Db() *database.Database
	Cache() PermissionCache
	IsBotAdmin(uint64) bool
	GetGuild(uint64) (guild.Guild, error)
	GetChannel(uint64) (channel.Channel, error)
	GetGuildMember(guildId, userId uint64) (member.Member, error)
	GetGuildRoles(uint64) ([]guild.Role, error)
	GetGuildOwner(uint64) (uint64, error)
}

func GetPermissionLevel(retriever Retriever, member member.Member, guildId uint64) (permLevel PermissionLevel, returnedError error) {
	// Check user ID in cache
	if cached, err := retriever.Cache().GetCachedPermissionLevel(guildId, member.User.Id); err == nil {
		return cached, nil
	} else if err != ErrNotCached {
		return Everyone, err
	}

	// Check if the user is a bot admin user
	if retriever.IsBotAdmin(member.User.Id) {
		return Admin, nil
	}

	// Don't recache if already cached (for now?)
	defer func() {
		if returnedError == nil {
			returnedError = retriever.Cache().SetCachedPermissionLevel(guildId, member.User.Id, permLevel)
		}
	}()

	// Check if user is guild owner
	if guildOwner, err := retriever.GetGuildOwner(guildId); err == nil {
		if member.User.Id == guildOwner {
			return Admin, nil
		}
	} else {
		return Everyone, err
	}

	// Check user perms for admin
	if adminUser, err := retriever.Db().Permissions.IsAdmin(guildId, member.User.Id); err == nil {
		if adminUser {
			return Admin, nil
		}
	} else {
		return Everyone, err
	}

	// Check roles from DB
	adminRoles, err := retriever.Db().RolePermissions.GetAdminRoles(guildId)
	if err != nil {
		return Everyone, err
	}

	for _, adminRoleId := range adminRoles {
		if member.HasRole(adminRoleId) {
			return Admin, nil
		}
	}

	// Check if user has Administrator permission
	hasAdminPermission := HasPermissions(retriever, guildId, member.User.Id, permission.Administrator)
	if hasAdminPermission {
		return Admin, nil
	}

	// Check user perms for support
	if isSupport, err := retriever.Db().Permissions.IsSupport(guildId, member.User.Id); err == nil {
		if isSupport {
			return Support, nil
		}
	} else {
		return Everyone, err
	}

	// Check if user is a member of a support team
	if isSupport, err := retriever.Db().SupportTeamMembers.IsSupport(guildId, member.User.Id); err == nil {
		if isSupport {
			return Support, nil
		}
	} else {
		return Everyone, err
	}

	// Check DB for support roles
	supportRoles, err := retriever.Db().RolePermissions.GetSupportRoles(guildId)
	if err != nil {
		return Everyone, err
	}

	for _, supportRoleId := range supportRoles {
		if member.HasRole(supportRoleId) {
			return Support, nil
		}
	}

	// Check if user has a role assigned to a support team
	if isSupport, err := retriever.Db().SupportTeamRoles.IsSupportAny(guildId, member.Roles); err == nil {
		if isSupport {
			return Support, nil
		}
	} else {
		return Everyone, err
	}

	return Everyone, nil
}
