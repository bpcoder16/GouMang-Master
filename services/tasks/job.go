package tasks

import (
	"context"
	"errors"
	"goumang-master/db"
	"goumang-master/global"
	"strconv"
	"time"

	"github.com/bpcoder16/Chestnut/v2/contrib/cron"
	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func getJobDefinition(ctx context.Context, dbTask db.GMTask) (jobDefinition gocron.JobDefinition, err error) {
	switch dbTask.Type {
	case db.TypeCron:
		if err = IsValidCrontabExpression(ctx, dbTask.Expression); err != nil {
			return
		}
		jobDefinition = gocron.CronJob(
			dbTask.Expression,
			true,
		)
	case db.TypeDuration:
		if durationMillisecond, errD := IsValidDurationExpression(ctx, dbTask.Expression); errD != nil {
			err = errD
		} else {
			jobDefinition = gocron.DurationJob(durationMillisecond)
		}
	case db.TypeDurationRandom:
		if minDurationMillisecond, maxDurationMillisecond, errD := IsValidDurationRandomExpression(ctx, dbTask.Expression); errD != nil {
			err = errD
		} else {
			jobDefinition = gocron.DurationRandomJob(minDurationMillisecond, maxDurationMillisecond)
		}
	case db.TypeOneTimeJobStartImmediately:
		jobDefinition = gocron.OneTimeJob(
			gocron.OneTimeJobStartImmediately(),
		)
	case db.TypeOneTimeJobStartDateTimes:
		if startAtList, errD := IsValidOneTimeJobStartDateTimesExpression(ctx, dbTask.Expression); errD != nil {
			err = errD
		} else {
			jobDefinition = gocron.OneTimeJob(
				gocron.OneTimeJobStartDateTimes(startAtList...),
			)
		}
	default:
		logit.Context(ctx).WarnW("getJobDefinition.Err", "["+strconv.FormatUint(uint64(dbTask.Type), 10)+"] "+notSupportedTaskTypeErr.Error())
		err = notSupportedTaskTypeErr
	}
	return
}

func getJobOptionList(ctx context.Context, taskUUID uuid.UUID, dbTask db.GMTask) []gocron.JobOption {
	return []gocron.JobOption{
		gocron.WithContext(ctx),
		gocron.WithName(dbTask.Title + taskNameDelimiter + dbTask.SHA256),
		gocron.WithSingletonMode(gocron.LimitModeReschedule), // LimitModeReschedule 重新调度模式(无法执行跳过等待下次周期尝试) LimitModeWait 等待模式(放入队列等待执行)
		gocron.WithIdentifier(taskUUID),
		gocron.WithEventListeners(
			// Job 执行前执行
			gocron.BeforeJobRuns(beforeJobRunsFunc(ctx)),

			// Job 执行前执行，抛出 error 可打断执行
			gocron.BeforeJobRunsSkipIfBeforeFuncErrors(beforeJobRunsSkipIfBeforeFuncErrorsFunc(ctx)),

			// Job 执行后执行
			gocron.AfterJobRuns(afterJobRunsFunc(ctx)),

			// Job 执行后执行，接收 Error
			gocron.AfterJobRunsWithError(afterJobRunsWithErrorFunc(ctx)),

			// Job 执行后执行，接收 Panic
			gocron.AfterJobRunsWithPanic(afterJobRunsWithPanicFunc(ctx)),

			//gocron.AfterLockError(func(jobID uuid.UUID, jobName string, err error) {
			//	logit.DebugW("AfterLockError.jobID", jobID, "uuid.UUID", jobName, "error", err)
			//}),
		),
	}
}

func CreateJob(ctx context.Context, dbTask db.GMTask) (job gocron.Job, err error) {
	var taskUUID uuid.UUID
	taskUUID, err = isValidTaskUUID(ctx, dbTask.UUID)
	if err != nil {
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
	if err != nil {
		return
	}

	var nextRunTime time.Time
	if nextRunTime, err = job.NextRun(); err != nil {
		return
	}

	dbTask.NextRunTime = nextRunTime.UnixNano() / 1e6
	dbTask.UpdatedAt = time.Now().UnixNano() / 1e6
	err = global.DefaultDB.WithContext(ctx).Save(&dbTask).Error
	if err == nil {
		logit.Context(ctx).DebugW("Cron.CreateJob", "Created: "+dbTask.Title)
	}
	return
}

func UpdateJob(ctx context.Context, dbTask db.GMTask) (job gocron.Job, err error) {
	var taskUUID uuid.UUID
	taskUUID, err = isValidTaskUUID(ctx, dbTask.UUID)
	if err != nil {
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

	job, err = cron.Update(
		taskUUID,
		jobDefinition,
		task,
		getJobOptionList(ctx, taskUUID, dbTask)...,
	)
	if err == nil {
		logit.Context(ctx).DebugW("Cron.UpdateJob", "Updated: "+dbTask.Title)
	}
	return
}

func RemoveJobForJob(ctx context.Context, job gocron.Job) (err error) {
	err = cron.RemoveJob(job.ID())
	logit.Context(ctx).DebugW("Cron.reloadTaskListTask", "Removed: "+job.Name())
	return
}

func RemoveJobForDBTask(ctx context.Context, dbTask db.GMTask) (err error) {
	var taskUUID uuid.UUID
	taskUUID, err = isValidTaskUUID(ctx, dbTask.UUID)
	if err != nil {
		return
	}
	err = cron.RemoveJob(taskUUID)
	logit.Context(ctx).DebugW("Cron.reloadTaskListTask", "Removed: "+dbTask.Title)
	return
}

func GetJob(jobID string) (job gocron.Job, err error) {
	jobList := cron.Jobs()
	for _, jobTmp := range jobList {
		if jobID == jobTmp.ID().String() {
			job = jobTmp
			return
		}
	}
	err = errors.New("job not found")
	return
}
