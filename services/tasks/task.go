package tasks

import (
	"context"
	"errors"
	"goumang-master/db"
	"goumang-master/global"
	"strconv"
	"strings"
	"time"

	"github.com/bpcoder16/Chestnut/v2/contrib/cron"
	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/go-co-op/gocron/v2"
)

const (
	taskNameDelimiter = "X_|_X"
)

func getTask(ctx context.Context, dbTask db.GMTask) (task gocron.Task, err error) {
	switch dbTask.Method {
	case db.MethodTest:
		task = testTask(dbTask)
	case db.MethodReloadTaskList:
		task = reloadTaskListTask(dbTask)
	default:
		logit.Context(ctx).WarnW("getTask.Err", "["+strconv.FormatUint(uint64(dbTask.Method), 10)+"] "+notSupportedTaskMethodErr.Error())
		err = notSupportedTaskMethodErr
	}
	return
}

func testTask(masterTask db.GMTask) (task gocron.Task) {
	task = gocron.NewTask(func(ctx context.Context) {
		logit.Context(ctx).DebugW("Cron.testTask", masterTask.Title+".Run")
	})
	return
}

func loadTaskListTask(ctx context.Context, dbTaskList []db.GMTask, exceptUUID string) {
	jobList := cron.Jobs()
	if len(jobList) == 0 && len(dbTaskList) == 0 {
		return
	}

	dbTaskMap := make(map[string]db.GMTask, len(dbTaskList))
	for _, dbTaskTmp := range dbTaskList {
		dbTaskMap[dbTaskTmp.UUID] = dbTaskTmp
	}
	dbTaskList = nil

	jobMap := make(map[string]gocron.Job, len(jobList))
	for _, job := range jobList {
		if job.ID().String() == exceptUUID {
			continue
		}
		// 移除无效的任务
		if _, isExist := dbTaskMap[job.ID().String()]; !isExist {
			err := RemoveJobForJob(ctx, job)
			logit.Context(ctx).DebugW("Cron.reloadTaskListTask", "Removed: "+job.Name())
			if err != nil {
				logit.Context(ctx).WarnW("Cron.reloadTaskListTask.RemoveJob.Err", err.Error())
			}
		} else {
			jobMap[job.ID().String()] = job
		}
	}
	jobList = nil

	for _, dbTask := range dbTaskMap {
		// 增加新任务
		if jobTmp, isExist := jobMap[dbTask.UUID]; !isExist {
			_, err := CreateJob(ctx, dbTask)
			if err = cancelErrJob(ctx, err, dbTask); err != nil {
				logit.Context(ctx).WarnW("Cron.reloadTaskListTask.cancelErrJob.Err", err.Error())
			}
		} else {
			if dbTask.Type == db.TypeOneTimeJobStartDateTimes {
				_, err := getJobDefinition(ctx, dbTask)
				if err = cancelErrJob(ctx, err, dbTask); err != nil {
					logit.Context(ctx).WarnW("Cron.reloadTaskListTask.cancelErrJob.Err", err.Error())
					continue
				}
			}
			// 判断是否需要更新
			nameList := strings.Split(jobTmp.Name(), taskNameDelimiter)
			if len(nameList) > 0 && nameList[len(nameList)-1] == dbTask.SHA256 {
				logit.Context(ctx).DebugW("Cron.reloadTaskListTask", "Unchanged: "+dbTask.Title)
				continue
			}

			_, err := UpdateJob(ctx, dbTask)
			if err = cancelErrJob(ctx, err, dbTask); err != nil {
				logit.Context(ctx).WarnW("Cron.reloadTaskListTask.cancelErrJob.Err", err.Error())
			}
		}
	}

	logit.Context(ctx).DebugW("Cron.reloadTaskListTask.Jobs.Final.Count", len(cron.Jobs()))
}

func reloadTaskListTask(masterTask db.GMTask) (task gocron.Task) {
	task = gocron.NewTask(func(ctx context.Context) {
		var dbTaskList []db.GMTask
		if err := global.DefaultDB.WithContext(ctx).Where("status = ? and id != ?", db.StatusEnabled, masterTask.ID).
			Order("id asc").Find(&dbTaskList).Error; err != nil {
			logit.Context(ctx).ErrorW("Cron.reloadTaskListTask.dbTaskList.Err", err.Error())
			return
		}
		loadTaskListTask(ctx, dbTaskList, masterTask.UUID)
	})
	return
}

func cancelErrJob(ctx context.Context, cancelErr error, task db.GMTask) error {
	cancelFunc := func(status int8) error {
		task.Status = status
		task.ErrorMessage = cancelErr.Error()
		task.UpdatedAt = uint64(time.Now().Unix())
		if err := global.DefaultDB.WithContext(ctx).Save(&task).Error; err != nil {
			return err
		}
		return RemoveJobForDBTask(ctx, task)
	}
	switch {
	case errors.Is(cancelErr, crontabExpressionErr),
		errors.Is(cancelErr, durationExpressionErr),
		errors.Is(cancelErr, durationRandomExpressionErr),
		errors.Is(cancelErr, oneTimeJobStartDateTimesExpressionErr),
		errors.Is(cancelErr, notSupportedTaskTypeErr),
		errors.Is(cancelErr, notSupportedTaskMethodErr),
		errors.Is(cancelErr, dbTaskUUIDInvalidErr):
		return cancelFunc(db.StatusConfigError)
	case errors.Is(cancelErr, oneTimeJobStartDateTimesExpressionExpired):
		return cancelFunc(db.StatusConfigExpired)
	default:
		return nil
	}
}
