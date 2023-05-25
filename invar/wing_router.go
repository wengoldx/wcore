// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2019/06/30   zhaixing       Add function from godfs
// -------------------------------------------------------------------

package invar

// RBAC role string
const (
	WRoleSuper     = "super-admin"   // Super admin, auto add
	WRoleAdmin     = "admin"         // Normal admin, same time as super admin permissions
	WRoleUser      = "user"          // Default normal user
	WRoleMComp     = "mall-comp"     // Mall composer account
	WRoleMDesigner = "mall-designer" // Mall designer account
	WRoleSComp     = "store-comp"    // Store composer account
	WRoleSMachine  = "store-machine" // Store machine account
	WRoleQKPartner = "qk-partner"    // QKS partner account
)

// RBAC role router keyword
const (
	WRGroupUser     = "user"
	WRGroupAdmin    = "admin"
	WRGroupComp     = "comp"
	WRGroupDesigner = "design"
	WRGroupMachine  = "mach"
	WRGroupPartner  = "part"
)

// Return role router key by given role, it maybe just return
// role string when not found from defined roles
func GetRouterKey(role string) string {
	switch role {
	case WRoleSuper, WRoleAdmin:
		return WRGroupAdmin
	case WRoleUser:
		return WRGroupUser
	case WRoleMComp, WRoleSComp:
		return WRGroupComp
	case WRoleMDesigner:
		return WRGroupDesigner
	case WRoleSMachine:
		return WRGroupMachine
	case WRoleQKPartner:
		return WRGroupPartner
	}
	return role
}
