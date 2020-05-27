package permission

import (
	"errors"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/permission"
)

func HasPermissionsChannel(retriever Retriever, guildId, userId, channelId uint64, permissions ...permission.Permission) bool {
	sum, err := GetEffectivePermissionsChannel(retriever, guildId, userId, channelId)
	if err != nil {
		return false
	}

	if permission.HasPermissionRaw(sum, permission.Administrator) {
		return true
	}

	hasPermission := true

	for _, perm := range permissions {
		if !permission.HasPermissionRaw(sum, perm) {
			hasPermission = false
			break
		}
	}

	return hasPermission
}

func HasPermissions(retriever Retriever, guildId, userId uint64, permissions ...permission.Permission) bool {
	sum, err := GetEffectivePermissions(retriever, guildId, userId)
	if err != nil {
		return false
	}

	if permission.HasPermissionRaw(sum, permission.Administrator) {
		return true
	}

	hasPermission := true

	for _, perm := range permissions {
		if !permission.HasPermissionRaw(sum, perm) {
			hasPermission = false
			break
		}
	}

	return hasPermission
}

func GetAllPermissionsChannel(retriever Retriever, guildId, userId, channelId uint64) []permission.Permission {
	permissions := make([]permission.Permission, 0)

	sum, err := GetEffectivePermissionsChannel(retriever, guildId, userId, channelId)
	if err != nil {
		return permissions
	}

	for _, perm := range permission.AllPermissions {
		if permission.HasPermissionRaw(sum, perm) {
			permissions = append(permissions, perm)
		}
	}

	return permissions
}

func GetAllPermissions(retriever Retriever, guildId, userId uint64) []permission.Permission {
	permissions := make([]permission.Permission, 0)

	sum, err := GetEffectivePermissions(retriever, guildId, userId)
	if err != nil {
		return permissions
	}

	for _, perm := range permission.AllPermissions {
		if permission.HasPermissionRaw(sum, perm) {
			permissions = append(permissions, perm)
		}
	}

	return permissions
}

func GetEffectivePermissionsChannel(retriever Retriever, guildId, userId, channelId uint64) (int, error) {
	permissions, err := GetBasePermissions(retriever, guildId)
	if err != nil {
		return 0, err
	}

	permissions, err = GetGuildTotalRolePermissions(retriever, guildId, userId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelBasePermissions(retriever, guildId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelTotalRolePermissions(retriever, guildId, userId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	permissions, err = GetChannelMemberPermissions(retriever, userId, channelId, permissions)
	if err != nil {
		return 0, err
	}

	return permissions, nil
}

func GetEffectivePermissions(retriever Retriever, guildId, userId uint64) (int, error) {
	permissions, err := GetBasePermissions(retriever, guildId)
	if err != nil {
		return 0, err
	}

	permissions, err = GetGuildTotalRolePermissions(retriever, guildId, userId, permissions)
	if err != nil {
		return 0, err
	}

	return permissions, nil
}

func GetChannelMemberPermissions(retriever Retriever, userId, channelId uint64, initialPermissions int) (int, error) {
	ch, err := retriever.GetChannel(channelId)
	if err != nil {
		return 0, err
	}

	for _, overwrite := range ch.PermissionOverwrites {
		if overwrite.Type == channel.PermissionTypeMember && overwrite.Id == userId {
			initialPermissions &= ^overwrite.Deny
			initialPermissions |= overwrite.Allow
		}
	}

	return initialPermissions, nil
}

func GetChannelTotalRolePermissions(retriever Retriever, guildId, userId, channelId uint64, initialPermissions int) (int, error) {
	member, err := retriever.GetGuildMember(guildId, userId)
	if err != nil {
		return 0, err
	}

	roles, err := retriever.GetGuildRoles(guildId)
	if err != nil {
		return 0, err
	}

	ch, err := retriever.GetChannel(channelId)
	if err != nil {
		return 0, err
	}

	allow, deny := 0, 0

	for _, memberRole := range member.Roles {
		for _, role := range roles {
			if memberRole == role.Id {
				for _, overwrite := range ch.PermissionOverwrites {
					if overwrite.Type == channel.PermissionTypeRole && overwrite.Id == role.Id {
						allow |= overwrite.Allow
						deny |= overwrite.Deny
						break
					}
				}
			}
		}
	}

	initialPermissions &= ^deny
	initialPermissions |= allow

	return initialPermissions, nil
}

func GetChannelBasePermissions(retriever Retriever, guildId, channelId uint64, initialPermissions int) (int, error) {
	roles, err := retriever.GetGuildRoles(guildId)
	if err != nil {
		return 0, err
	}

	var publicRole *guild.Role
	for _, role := range roles {
		if role.Id == guildId {
			publicRole = &role
			break
		}
	}

	if publicRole == nil {
		return 0, errors.New("couldn't find public role")
	}

	ch, err := retriever.GetChannel(channelId)
	if err != nil {
		return 0, err
	}

	for _, overwrite := range ch.PermissionOverwrites {
		if overwrite.Type == channel.PermissionTypeRole && overwrite.Id == publicRole.Id {
			initialPermissions &= ^overwrite.Deny
			initialPermissions |= overwrite.Allow
			break
		}
	}

	return initialPermissions, nil
}

func GetGuildTotalRolePermissions(retriever Retriever, guildId, userId uint64, initialPermissions int) (int, error) {
	member, err := retriever.GetGuildMember(guildId, userId)
	if err != nil {
		return 0, err
	}

	roles, err := retriever.GetGuildRoles(guildId)
	if err != nil {
		return 0, err
	}

	for _, memberRole := range member.Roles {
		for _, role := range roles {
			if memberRole == role.Id {
				initialPermissions |= role.Permissions
			}
		}
	}

	return initialPermissions, nil
}

func GetBasePermissions(retriever Retriever, guildId uint64) (int, error) {
	roles, err := retriever.GetGuildRoles(guildId)
	if err != nil {
		return 0, err
	}

	var publicRole *guild.Role
	for _, role := range roles {
		if role.Id == guildId {
			publicRole = &role
			break
		}
	}

	if publicRole == nil {
		return 0, errors.New("couldn't find public role")
	}

	return publicRole.Permissions, nil
}

