package svc

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/zeromicro/go-zero/rest"

	"neekity.com/go-admin/api/internal/config"
	"neekity.com/go-admin/api/internal/middleware"
	"neekity.com/go-admin/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config            config.Config
	AccessLog         rest.Middleware
	Auth              rest.Middleware
	CasbinRuleEnforce *casbin.Enforcer
	UserModel         model.UserModel
	RoleModel         model.RoleModel
	PermissionModel   model.PermissionModel
	MenuModel         model.MenuModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	gdb, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	a, err := gormadapter.NewAdapterByDBUseTableName(gdb, "", "casbin_rule")
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer("/data/resource/rbac_model.conf", a)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:            c,
		AccessLog:         middleware.NewAccessLogMiddleware().Handle,
		Auth:              middleware.NewAuthMiddleware(c.JwtAuth.AccessSecret, e).Handle,
		CasbinRuleEnforce: e,
		UserModel:         model.NewUserModel(gdb),
		RoleModel:         model.NewRoleModel(gdb),
		PermissionModel:   model.NewPermissionModel(gdb),
		MenuModel:         model.NewMenuModel(gdb),
	}
}
