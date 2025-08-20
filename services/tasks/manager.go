package tasks

import (
	"context"
	"goumang-master/db"
	"goumang-master/global"

	"github.com/bpcoder16/Chestnut/v2/contrib/cron"
	"github.com/bpcoder16/Chestnut/v2/logit"
)

func InitJobs(ctx context.Context) {
	var dbTaskList []db.GMTask
	if err := global.DefaultDB.WithContext(ctx).Where("status = ?", db.StatusEnabled).
		Order("id asc").Find(&dbTaskList).Error; err != nil {
		panic("Cron.InitJobs.获取任务列表异常, Err: " + err.Error())
	}

	for _, dbTaskItem := range dbTaskList {
		_, err := CreateJob(ctx, dbTaskItem)
		if err != nil {
			panic("Cron.InitJobs.创建任务失败, Err: " + err.Error())
		}
	}

	logit.Context(ctx).InfoW("Cron.InitJobs.Count", len(cron.Jobs()))
}
