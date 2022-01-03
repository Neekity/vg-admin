package model

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
	"neekity.com/go-admin/common/constant"
	"neekity.com/go-admin/common/helper"
)

type (
	RoleModel interface {
		FindOne(id uint, enforce *casbin.Enforcer) (*RoleInfo, error)
		List(name string) ([]Role, error)
		UpdateOrCreate(id uint, RoleInfo Role, permissionIds []int) (*Role, error)
		Delete(id uint) error
	}

	defaultRoleModel struct {
		GormDB *gorm.DB
		table  string
	}

	Role struct {
		ID         uint   `gorm:"column:id;type:uint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
		Name       string `gorm:"type:string;comment:角色名称;size:64;not null;" json:"name"`
		CasbinRole string `gorm:"type:string;uniqueIndex;comment:角色key;size:64;not null;" json:"casbin_role"`
		CreatedAt  helper.MyTime
		UpdatedAt  helper.MyTime
	}

	RoleInfo struct {
		Role
		Permissions []Permission `json:"permissions"`
	}
)

func (m *defaultRoleModel) Delete(id uint) error {
	return m.GormDB.Delete(&Role{ID: id}).Error
}

func (m *defaultRoleModel) UpdateOrCreate(id uint, RoleInfo Role, permissionIds []int) (*Role, error) {
	var role Role
	var permissions []Permission

	err := m.GormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Role{}).Where("id = ?", id).Assign(RoleInfo).FirstOrCreate(&role).Error; err != nil {
			return err
		}
		if err := m.GormDB.Where("ptype = ?", constant.CasbinPolicy).
			Where("v0 = ?", role.CasbinRole).Delete(CasbinRule{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&Permission{}).Find(&permissions, permissionIds).Error; err != nil {
			return err
		}
		var casbinModels []CasbinRule
		for _, permission := range permissions {
			casbinModels = append(casbinModels, CasbinRule{
				Ptype: constant.CasbinPolicy,
				V0:    role.CasbinRole,
				V1:    permission.CasbinPermission,
				V2:    permission.CasbinPermissionType,
			})
		}
		if err := tx.Create(&casbinModels).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (m *defaultRoleModel) List(name string) ([]Role, error) {
	var roles []Role
	query := m.GormDB.Table(m.table)
	if name != "" {
		query.Scopes(helper.QueryKey("%"+name+"%", "name", "like"))
	}
	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (m *defaultRoleModel) FindOne(id uint, enforce *casbin.Enforcer) (*RoleInfo, error) {
	var role Role
	var permissions []Permission
	if err := m.GormDB.First(&role, id).Error; err != nil {
		return nil, err
	}
	if err := enforce.LoadPolicy(); err != nil {
		return nil, err
	}
	policies := enforce.GetFilteredNamedPolicy(constant.CasbinPolicy, 0, role.CasbinRole)
	for _, policy := range policies {
		var permission Permission
		err := m.GormDB.Model(&Permission{}).Where("casbin_permission = ?", policy[1]).
			Where("casbin_permission_type = ?", policy[2]).First(&permission).Error
		if err != nil {
			if err == GormErrNotFound {
				continue
			}
			return nil, err

		}
		permissions = append(permissions, permission)
	}

	roleInfo := RoleInfo{
		role,
		permissions,
	}

	return &roleInfo, nil
}

func NewRoleModel(gdb *gorm.DB) RoleModel {
	gdb.Set("gorm:table_options", "ENGINE=InnoDB COMMENT '角色表'  charset = utf8mb4;").
		AutoMigrate(&Role{})
	return &defaultRoleModel{
		GormDB: gdb,
		table:  "`role`",
	}
}
