package tasks

import (
	"context"
	"fmt"
	"goumang-master/db"
	"goumang-master/global"
	"strconv"
	"time"

	"github.com/bpcoder16/Chestnut/v2/appconfig/env"
	"github.com/bpcoder16/Chestnut/v2/contrib/cron"
	"github.com/bpcoder16/Chestnut/v2/logit"
	robfigCron "github.com/robfig/cron/v3"
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

func IsValidCrontabExpression(expression string) (err error) {
	p := robfigCron.NewParser(robfigCron.SecondOptional | robfigCron.Minute | robfigCron.Hour | robfigCron.Dom | robfigCron.Month | robfigCron.Dow | robfigCron.Descriptor)
	withLocation := fmt.Sprintf("CRON_TZ=%s %s", env.TimeLocation().String(), expression)
	_, err = p.Parse(withLocation)
	return
}

func IsValidDurationExpression(expression string) (durationMillisecond time.Duration, err error) {
	var durationMillisecondInt int
	durationMillisecondInt, err = strconv.Atoi(expression)
	if err != nil {
		return
	}

	durationMillisecond = time.Duration(durationMillisecondInt) * time.Millisecond
	return
}
