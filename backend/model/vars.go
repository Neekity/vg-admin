package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var ErrNotFound = sqlx.ErrNotFound
var GormErrNotFound = gorm.ErrRecordNotFound
