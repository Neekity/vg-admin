package model

import (
	"gorm.io/gorm"
	"neekity.com/go-admin/common/helper"
)

type (
	PermissionModel interface {
		FindOne(id uint) (*Permission, error)
		List(name string) ([]Permission, error)
		UpdateOrCreate(id uint, permissionInfo Permission) (*Permission, error)
		Delete(id uint) error
	}

	defaultPermissionModel struct {
		GormDB *gorm.DB
		table  string
	}

	Permission struct {
		ID                   uint   `gorm:"column:id;type:uint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
		Name                 string `gorm:"type:string;comment:角色名称;size:64;not null;" json:"name"`
		CasbinPermission     string `gorm:"type:string;uniqueIndex:udx_casbin_policy;comment:casbin权限key;size:64;not null;" json:"casbin_permission"`
		CasbinPermissionType string `gorm:"type:string;uniqueIndex:udx_casbin_policy;comment:casbin权限类型;size:64;not null;" json:"casbin_permission_type"`
		Route                string `gorm:"type:string;comment:路由;size:64;not null;" json:"route"`
		CreatedAt            helper.MyTime
		UpdatedAt            helper.MyTime
	}
)

func (m *defaultPermissionModel) Delete(id uint) error {
	return m.GormDB.Delete(&Permission{ID: id}).Error
}

func (m *defaultPermissionModel) UpdateOrCreate(id uint, PermissionInfo Permission) (*Permission, error) {
	var permission Permission

	if err := m.GormDB.Model(&Permission{}).Where("id = ?", id).Assign(PermissionInfo).FirstOrCreate(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

func (m *defaultPermissionModel) List(name string) ([]Permission, error) {
	var permissions []Permission
	query := m.GormDB.Table(m.table)
	if name != "" {
		query.Scopes(helper.QueryKey("%"+name+"%", "name", "like"))
	}
	if err := query.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (m *defaultPermissionModel) FindOne(id uint) (*Permission, error) {
	var permission Permission
	if err := m.GormDB.First(&permission, id).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

func NewPermissionModel(gdb *gorm.DB) PermissionModel {
	gdb.Set("gorm:table_options", "ENGINE=InnoDB COMMENT '权限表'  charset = utf8mb4;").
		AutoMigrate(&Permission{})
	return &defaultPermissionModel{
		GormDB: gdb,
		table:  "`permission`",
	}
}
