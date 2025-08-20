package tasks

import (
	"context"
	"errors"
	"fmt"
	"goumang-master/db"
	"strconv"
	"strings"
	"time"

	"github.com/bpcoder16/Chestnut/v2/appconfig/env"
	"github.com/bpcoder16/Chestnut/v2/contrib/cron"
	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	robfigCron "github.com/robfig/cron/v3"
)

func IsValidCrontabExpression(expression string) (err error) {
	p := robfigCron.NewParser(robfigCron.SecondOptional | robfigCron.Minute | robfigCron.Hour | robfigCron.Dom | robfigCron.Month | robfigCron.Dow | robfigCron.Descriptor)
	withLocation := fmt.Sprintf("CRON_TZ=%s %s", env.TimeLocation().String(), expression)
	_, err = p.Parse(withLocation)
	if err != nil {
		err = errors.New("[" + expression + "]" + err.Error())
	}
	return
}

func IsValidDurationExpression(expression string) (durationMillisecond time.Duration, err error) {
	var durationMillisecondInt int
	durationMillisecondInt, err = strconv.Atoi(expression)
	if err != nil {
		err = errors.New("[" + expression + "]" + err.Error())
		return
	}

	durationMillisecond = time.Duration(durationMillisecondInt) * time.Millisecond
	return
}

func IsValidDurationRandomExpression(expression string) (minDurationMillisecond, maxDurationMillisecond time.Duration, err error) {
	expressionList := strings.Split(expression, ",")
	if len(expressionList) != 2 {
		err = errors.New("[" + expression + "]invalid expression")
		return
	}
	var minDurationMillisecondInt, maxDurationMillisecondInt int
	minDurationMillisecondInt, err = strconv.Atoi(expressionList[0])
	if err != nil {
		err = errors.New("[" + expression + "]" + err.Error())
		return
	}
	maxDurationMillisecondInt, err = strconv.Atoi(expressionList[1])
	if err != nil {
		err = errors.New("[" + expression + "]" + err.Error())
		return
	}
	if minDurationMillisecondInt >= maxDurationMillisecondInt || minDurationMillisecondInt <= 0 || maxDurationMillisecondInt <= 0 {
		err = errors.New("[" + expression + "]invalid expression")
		return
	}
	minDurationMillisecond = time.Duration(minDurationMillisecondInt) * time.Millisecond
	maxDurationMillisecond = time.Duration(maxDurationMillisecondInt) * time.Millisecond
	return
}

func IsValidOneTimeJobStartDateTimesExpression(expression string) (timeList []time.Time, err error) {
	expressionList := strings.Split(expression, ",")
	if len(expressionList) == 0 {
		err = errors.New("[" + expression + "]invalid expression")
		return
	}

	timeList = make([]time.Time, 0, len(expressionList))
	for _, timeStr := range expressionList {
		startAt, errT := time.ParseInLocation(time.DateTime, timeStr, env.TimeLocation())
		if errT != nil {
			err = errors.New("[" + expression + "]" + errT.Error())
			return
		}
		timeList = append(timeList, startAt)
	}

	maxTime := timeList[0]
	for _, t := range timeList[1:] {
		if t.After(maxTime) {
			maxTime = t
		}
	}
	if maxTime.Before(time.Now()) {
		err = errors.New("[" + expression + "]max time is too early")
	}

	return
}

func getJobDefinition(_ context.Context, dbTask db.GMTask) (jobDefinition gocron.JobDefinition, err error) {
	switch dbTask.Type {
	case db.TypeCron:
		if err = IsValidCrontabExpression(dbTask.Expression); err != nil {
			return
		}
		jobDefinition = gocron.CronJob(
			dbTask.Expression,
			true,
		)
	case db.TypeDuration:
		if durationMillisecond, errD := IsValidDurationExpression(dbTask.Expression); errD != nil {
			err = errD
		} else {
			jobDefinition = gocron.DurationJob(durationMillisecond)
		}
	case db.TypeDurationRandom:
		if minDurationMillisecond, maxDurationMillisecond, errD := IsValidDurationRandomExpression(dbTask.Expression); errD != nil {
			err = errD
		} else {
			jobDefinition = gocron.DurationRandomJob(minDurationMillisecond, maxDurationMillisecond)
		}
	case db.TypeOneTimeJobStartImmediately:
		jobDefinition = gocron.OneTimeJob(
			gocron.OneTimeJobStartImmediately(),
		)
	case db.TypeOneTimeJobStartDateTimes:
		if startAtList, errD := IsValidOneTimeJobStartDateTimesExpression(dbTask.Expression); errD != nil {
			err = errD
		} else {
			jobDefinition = gocron.OneTimeJob(
				gocron.OneTimeJobStartDateTimes(startAtList...),
			)
		}
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
