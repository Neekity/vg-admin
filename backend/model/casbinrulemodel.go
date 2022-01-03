package model

import (
	"gorm.io/gorm"
)

type (
	CasbinRuleModel interface {
	}

	defaultCasbinRuleModel struct {
		GormDB *gorm.DB
		table  string
	}

	CasbinRule struct {
		ID    uint   `json:"id" gorm:"primaryKey;column:id;type:int(10) unsigned auto_increment"`
		Ptype string `json:"ptype" gorm:"column:ptype;type:varchar(255)"`
		V0    string `json:"v0" gorm:"column:v0;type:varchar(255)"`
		V1    string `json:"v1" gorm:"column:v1;type:varchar(255)"`
		V2    string `json:"v2" gorm:"column:v2;type:varchar(255)"`
		V3    string `json:"v3" gorm:"column:v3;type:varchar(255)"`
		V4    string `json:"v4" gorm:"column:v4;type:varchar(255)"`
		V5    string `json:"v5" gorm:"column:v5;type:varchar(255)"`
	}
)
