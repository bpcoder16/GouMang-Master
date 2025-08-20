package bootstrap

import (
	"context"
	"goumang-master/global"
	"goumang-master/services/tasks"
	"path"

	"github.com/bpcoder16/Chestnut/v2/appconfig"
	"github.com/bpcoder16/Chestnut/v2/appconfig/env"
	"github.com/bpcoder16/Chestnut/v2/bootstrap"
	"github.com/bpcoder16/Chestnut/v2/core/utils"
	"github.com/bpcoder16/Chestnut/v2/default/mysql"
	"github.com/bpcoder16/Chestnut/v2/default/sqlite"
)

func MustInit(ctx context.Context, config *appconfig.AppConfig) {
	bootstrap.MustInit(ctx, config)

	loadAppBizConfig()

	switch global.AppBizConfig.GormDBDriver {
	case "mysql":
		global.DefaultDB = mysql.MasterDB()
	default:
		global.DefaultDB = sqlite.DefaultClient()
	}

	tasks.InitJobs(ctx)
}

func loadAppBizConfig() {
	err := utils.ParseFile(path.Join(env.ConfigDirPath(), "app-biz-config.yaml"), &global.AppBizConfig)
	if err != nil {
		panic("load app-biz-config err:" + err.Error())
	}
}
