package tasks

import (
	"context"
	"errors"
	"goumang-master/db"
	"strconv"

	"github.com/bpcoder16/Chestnut/v2/contrib/cron"
	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func getJobDefinition(_ context.Context, dbTask db.GMTask) (jobDefinition gocron.JobDefinition, err error) {
	switch dbTask.Type {
	case db.TypeCron:
		if err = IsValidCrontab(dbTask.Expression); err != nil {
			return
		}
		jobDefinition = gocron.CronJob(
			dbTask.Expression,
			true,
		)
	default:
		err = errors.New("invalid task, type: " + strconv.FormatUint(uint64(dbTask.Type), 10))
	}
	return
}

func getJobOptionList(ctx context.Context, taskUUID uuid.UUID, dbTask db.GMTask) []gocron.JobOption {
	return []gocron.JobOption{
		gocron.WithContext(ctx),
		gocron.WithName(dbTask.Title + taskNameDelimiter + dbTask.SHA256),
		gocron.WithSingletonMode(gocron.LimitModeReschedule), // LimitModeReschedule 重新调度模式(无法执行跳过等待下次周期尝试) LimitModeWait 等待模式(放入队列等待执行)
		gocron.WithIdentifier(taskUUID),
	}
}

func CreateJob(ctx context.Context, dbTask db.GMTask) (job gocron.Job, err error) {
	logit.Context(ctx).DebugW("Cron.CreateJob.dbTask", dbTask)

	var taskUUID uuid.UUID
	taskUUID, err = uuid.Parse(dbTask.UUID)
	if err != nil {
		err = errors.New("invalid.task.uuid")
		return
	}

	var jobDefinition gocron.JobDefinition
	jobDefinition, err = getJobDefinition(ctx, dbTask)
	if err != nil {
		return
	}
	var task gocron.Task
	task, err = getTask(ctx, dbTask)
	if err != nil {
		return
	}

	job, err = cron.NewJob(
		jobDefinition,
		task,
		getJobOptionList(ctx, taskUUID, dbTask)...,
	)
	if err == nil {
		logit.Context(ctx).DebugW("Cron.CreateJob", "Created:"+dbTask.Title)
	}
	return
}
