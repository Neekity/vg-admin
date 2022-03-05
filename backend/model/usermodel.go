package model

import (
	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"neekity.com/go-admin/common/constant"
	"neekity.com/go-admin/common/helper"
	"strconv"
)

type (
	UserModel interface {
		FindOne(id int) (*User, error)
		List(name string) ([]User, error)
		UpdateOrCreate(email string, userInfo User) (*User, error)
		Delete(id uint) error
		GetUserRoles(id uint, enforce *casbin.Enforcer) ([]Role, error)
		Assign(id uint, casbinRoles []string) error
		GetUserInfo(userId uint, enforce *casbin.Enforcer) (*UserInfo, error)
	}

	defaultUserModel struct {
		GormDB *gorm.DB
		table  string
	}

	User struct {
		ID        uint   `gorm:"column:id;type:uint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
		Name      string `gorm:"type:string;comment:姓名;size:64;not null;" json:"name"`
		Email     string `gorm:"type:string;uniqueIndex;comment:邮箱;size:64;not null;" json:"email"`
		Password  string `gorm:"column:password;->:false;<-:update;<-:create" json:"password,omitempty"` // 密码
		Mobile    string `gorm:"type:string;comment:手机号;size:16;not null;" json:"mobile"`
		Avatar    string `gorm:"type:string;comment:头像;size:256;not null;" json:"avatar"`
		CreatedAt helper.MyTime
		UpdatedAt helper.MyTime
		DeletedAt gorm.DeletedAt
	}

	UserInfo struct {
		Access []string `json:"access"`
		User
	}
)

func (m *defaultUserModel) GetUserInfo(userId uint, enforce *casbin.Enforcer) (*UserInfo, error) {
	var user User
	var menuPath []string
	var access []string
	if err := m.GormDB.First(&user, userId).Error; err != nil {
		return nil, err
	}
	if err := enforce.LoadPolicy(); err != nil {
		return nil, err
	}
	if err := m.GormDB.Model(&Menu{}).Where("path != ?", "").Pluck("path", &menuPath).Error; err != nil {
		return nil, err
	}
	for _, item := range menuPath {
		if flag, err := enforce.Enforce(strconv.Itoa(int(userId)), item, constant.CasbinPermissionRead); err != nil {
			return nil, err
		} else if flag {
			access = append(access, item)
		} else {
			logx.Error(item, flag, err)
		}
	}
	if isAdmin := enforce.HasGroupingPolicy(strconv.Itoa(int(userId)), constant.AdminRole); isAdmin == true {
		access = menuPath
	}
	userInfo := UserInfo{
		access,
		user,
	}
	return &userInfo, nil
}

func (m *defaultUserModel) Assign(id uint, casbinRoles []string) error {
	return m.GormDB.Transaction(func(tx *gorm.DB) error {
		userId := strconv.Itoa(int(id))
		if err := tx.Where(CasbinRule{
			Ptype: constant.CasbinDefaultRole,
			V0:    userId,
		}).Delete(CasbinRule{}).Error; err != nil {
			return err
		}
		var models []CasbinRule
		for _, role := range casbinRoles {
			models = append(models, CasbinRule{
				Ptype: constant.CasbinDefaultRole,
				V0:    userId,
				V1:    role,
			})
		}

		if err := tx.Create(&models).Error; err != nil {
			return err
		}

		return nil
	})
}

func (m *defaultUserModel) GetUserRoles(id uint, enforce *casbin.Enforcer) ([]Role, error) {
	var roles []Role
	if err := enforce.LoadPolicy(); err != nil {
		return nil, err
	}
	policies := enforce.GetFilteredNamedGroupingPolicy(constant.CasbinDefaultRole, 0, strconv.Itoa(int(id)))
	for _, policy := range policies {
		var role Role

		if err := m.GormDB.Model(&Role{}).Where("casbin_role = ?", policy[1]).First(&role).Error; err != nil {
			if err == GormErrNotFound {
				continue
			}
			return nil, err
		}

		roles = append(roles, role)
	}
	return roles, nil
}

func (m *defaultUserModel) Delete(id uint) error {
	return m.GormDB.Delete(&User{ID: id}).Error
}

func (m *defaultUserModel) UpdateOrCreate(email string, userInfo User) (*User, error) {
	var material User
	if err := m.GormDB.Model(&User{}).Where("email = ?", email).Assign(userInfo).FirstOrCreate(&material).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

func (m *defaultUserModel) List(name string) ([]User, error) {
	var users []User
	query := m.GormDB.Table(m.table)
	if name != "" {
		query.Scopes(helper.QueryKey("%"+name+"%", "name", "like"))
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (m *defaultUserModel) FindOne(id int) (*User, error) {
	var material User
	if err := m.GormDB.First(&material, id).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

func NewUserModel(gdb *gorm.DB) UserModel {
	gdb.Set("gorm:table_options", "ENGINE=InnoDB COMMENT '用户表'  charset = utf8mb4;").
		AutoMigrate(&User{})
	return &defaultUserModel{
		GormDB: gdb,
		table:  "`user`",
	}
}
