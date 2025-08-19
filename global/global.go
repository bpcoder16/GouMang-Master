package global

import (
	"goumang-master/types"

	"gorm.io/gorm"
)

var (
	AppBizConfig *types.AppBizConfig
	DefaultDB    *gorm.DB
)
