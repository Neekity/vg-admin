package model

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var ErrNotFound = sqlx.ErrNotFound
var GormErrNotFound = gorm.ErrRecordNotFound
