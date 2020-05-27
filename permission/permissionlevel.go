package permission

type PermissionLevel int

const(
	Everyone PermissionLevel = iota
	Support
	Admin
)

func (l PermissionLevel) Int() int {
	return int(l)
}
