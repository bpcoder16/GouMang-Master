package tasks

import (
	"context"
	"goumang-master/db"
	"goumang-master/global"

	"github.com/google/uuid"
)

func InitJobs(ctx context.Context) {
	var dbTaskList []db.GMTask
	if err := global.DefaultDB.WithContext(ctx).Where("status = ?", db.StatusEnabled).
		Order("id asc").Find(&dbTaskList).Error; err != nil {
		panic("Cron.InitJobs.dbTaskList.Err: " + err.Error())
	}

	loadTaskListTask(ctx, dbTaskList, "")
	_, _ = CreateJob(ctx, global.DefaultDB.WithContext(ctx), db.GMTask{
		UUID:   uuid.New().String(),
		SHA256: "Immediately",
		Title:  "initJobNextRunTime",
		Type:   db.TypeOneTimeJobStartImmediately,
		Method: db.MethodInitJobNextRunTime,
	})
}
