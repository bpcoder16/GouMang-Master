package tasks

import (
	"context"
	"goumang-master/db"
	"goumang-master/global"
)

func InitJobs(ctx context.Context) {
	var dbTaskList []db.GMTask
	if err := global.DefaultDB.WithContext(ctx).Where("status = ?", db.StatusEnabled).
		Order("id asc").Find(&dbTaskList).Error; err != nil {
		panic("Cron.InitJobs.dbTaskList.Err: " + err.Error())
	}

	loadTaskListTask(ctx, dbTaskList, "")
}
