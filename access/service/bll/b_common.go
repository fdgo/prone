package bll

import (
	"access/service/menu_service"
	"access/service/role_service"
	"access/service/user_service"
)

type Common struct {
	UserAPI *user_service.User `inject:""`
	RoleAPI *Role_service.Role `inject:""`
	MenuAPI *menu_service.Menu `inject:""`
}
