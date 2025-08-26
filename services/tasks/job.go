package tasks

import (
	"context"
	"errors"
	"goumang-master/db"
	"strconv"
	"time"

	"github.com/bpcoder16/Chestnut/v2/contrib/cron"
	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func updateDBTaskNextRunTime(_ context.Context, gormDB *gorm.DB, job gocron.Job, isSetEditableNo bool) (err error) {
	var nextRunTime time.Time
	if nextRunTime, err = job.NextRun(); err != nil {
		return
	}
	if nextRunTime.Before(time.Now().Add(-time.Second)) {
		err = errors.New("nextRunTime is in the past")
		return
	}

	updateValues := map[string]any{
		"next_run_time": nextRunTime.UnixNano() / int64(time.Millisecond),
	}
	if isSetEditableNo {
		updateValues["editable"] = db.EditableNo
	}
	return gormDB.Model(&db.GMTask{}).
		Where("uuid = ?", job.ID().String()).
		Updates(updateValues).Error
}

func CreateJob(ctx context.Context, gormDB *gorm.DB, dbTask db.GMTask) (job gocron.Job, err error) {
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
		err = createOrUpdateJobErr
		return
	}
	logit.Context(ctx).DebugW("Cron.CreateJob", "Created: "+dbTask.Title)

	_ = updateDBTaskNextRunTime(ctx, gormDB, job, false)
	return
}

func UpdateJob(ctx context.Context, gormDB *gorm.DB, dbTask db.GMTask) (job gocron.Job, err error) {
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

	if err != nil {
		err = createOrUpdateJobErr
		return
	}
	logit.Context(ctx).DebugW("Cron.UpdateJob", "Updated: "+dbTask.Title)

	_ = updateDBTaskNextRunTime(ctx, gormDB, job, false)

	return
}

func RemoveJobForJob(ctx context.Context, job gocron.Job) (err error) {
	err = cron.RemoveJob(job.ID())
	logit.Context(ctx).DebugW("Cron.reloadTaskListTask.RemoveJobForJob", "Removed: "+job.Name())
	return
}

func RemoveJobForDBTask(ctx context.Context, dbTask db.GMTask) (err error) {
	var taskUUID uuid.UUID
	taskUUID, err = isValidTaskUUID(ctx, dbTask.UUID)
	if err != nil {
		return
	}
	err = cron.RemoveJob(taskUUID)
	logit.Context(ctx).DebugW("Cron.reloadTaskListTask.RemoveJobForDBTask", "Removed: "+dbTask.Title)
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
