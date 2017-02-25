package auth

// Roles holds user roles
type Roles struct {
	RoleDatabaseManager bool `xml:"database-manager" json:"database-manager"`
	RoleImageManager    bool `xml:"image-manager" json:"image-manager"`
	RoleScriptManager   bool `xml:"script-manager" json:"script-manager"`
	RoleUserManager     bool `xml:"user-manager" json:"user-manager"`
}

const (
	// RoleDatabaseManager role to save, delete and get info from database
	RoleDatabaseManager = iota

	// RoleImageManager role to save, delete and get info from image
	RoleImageManager

	// RoleScriptManager role to save, delete and get info from script
	RoleScriptManager

	// RoleUserManager role to save, delete and get info of database users
	RoleUserManager
)

// CheckUserPermission checks if UserClaim has the role
func CheckUserPermission(user UserClaim, role int) bool {
	var result bool
	switch role {
	case RoleDatabaseManager:
		result = user.Roles.RoleDatabaseManager
	case RoleImageManager:
		result = user.Roles.RoleImageManager
	case RoleScriptManager:
		result = user.Roles.RoleScriptManager
	case RoleUserManager:
		result = user.Roles.RoleUserManager
	}
	return result
}
