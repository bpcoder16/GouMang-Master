package tasks

import (
	"context"
	"goumang-master/db"
	"strconv"

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
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				logit.DebugW("BeforeJobRuns.jobID", jobID, "jobName", jobName)
			}),
			gocron.BeforeJobRunsSkipIfBeforeFuncErrors(func(jobID uuid.UUID, jobName string) error {
				logit.DebugW("BeforeJobRunsSkipIfBeforeFuncErrors.jobID", jobID, "jobName", jobName)
				return nil
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				logit.DebugW("AfterJobRuns.jobID", jobID, "jobName", jobName)
			}),

			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				logit.ErrorW("AfterJobRunsWithError.jobID", jobID, "jobName", jobName, "error", err)
			}),
			gocron.AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {
				logit.ErrorW("AfterJobRunsWithPanic.jobID", jobID, "jobName", jobName, "recoverData", recoverData)
			}),
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
