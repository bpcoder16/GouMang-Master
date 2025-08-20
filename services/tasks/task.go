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
	"github.com/google/uuid"
)

const (
	taskNameDelimiter = "__A__"
)

func getTask(_ context.Context, dbTask db.GMTask) (task gocron.Task, err error) {
	switch dbTask.Method {
	case db.MethodTest:
		task = testTask(dbTask)
	case db.MethodReloadTaskList:
		task = reloadTaskListTask(dbTask)
	default:
		err = errors.New("invalid method, method: " + strconv.FormatUint(uint64(dbTask.Method), 10))
	}
	return
}

func testTask(masterTask db.GMTask) (task gocron.Task) {
	task = gocron.NewTask(func(ctx context.Context) {
		logit.Context(ctx).DebugW("Cron.testTask", masterTask.Title+".Run")
	})
	return
}

func reloadTaskListTask(masterTask db.GMTask) (task gocron.Task) {
	task = gocron.NewTask(func(ctx context.Context) {
		var dbTaskList []db.GMTask
		if err := global.DefaultDB.WithContext(ctx).Where("status = ? and id != ?", db.StatusEnabled, masterTask.ID).
			Order("id asc").Find(&dbTaskList).Error; err != nil {
			logit.Context(ctx).ErrorW("Cron.reloadTaskListTask.dbTaskList.Err", err.Error())
			return
		}

		logit.Context(ctx).DebugW("Cron.reloadTaskListTask.dbTaskList.Count", len(dbTaskList))
		jobList := cron.Jobs()
		if len(jobList) == 0 && len(jobList) == 0 {
			return
		}

		dbTaskMap := make(map[string]db.GMTask, len(dbTaskList))
		for _, dbTaskTmp := range dbTaskList {
			dbTaskMap[dbTaskTmp.UUID] = dbTaskTmp
		}
		dbTaskList = nil

		jobMap := make(map[string]gocron.Job, len(jobList))
		for _, job := range jobList {
			if job.ID().String() == masterTask.UUID {
				continue
			}
			// 移除无效的任务
			if _, isExist := dbTaskMap[job.ID().String()]; !isExist {
				err := cron.RemoveJob(job.ID())
				logit.Context(ctx).DebugW("Cron.reloadTaskListTask", "Removed:"+job.Name())
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
				if err != nil {
					logit.Context(ctx).WarnW("Cron.reloadTaskListTask.CreateJob.Err", err.Error())
				}
			} else {
				// 判断是否需要更新
				nameList := strings.Split(jobTmp.Name(), taskNameDelimiter)
				if len(nameList) > 0 && nameList[len(nameList)-1] == dbTask.SHA256 {
					logit.Context(ctx).DebugW("Cron.reloadTaskListTask", "Unchanged:"+jobTmp.Name())
					continue
				}

				var jobDefinitionNew gocron.JobDefinition
				var errNew error
				jobDefinitionNew, errNew = getJobDefinition(ctx, dbTask)
				if errNew != nil {
					if err := cancelErrJob(ctx, jobTmp.ID(), dbTask); err != nil {
						logit.Context(ctx).WarnW("Cron.reloadTaskListTask.cancelErrJob.Err", err.Error())
					}
					continue
				}
				var taskNew gocron.Task
				taskNew, errNew = getTask(ctx, dbTask)
				if errNew != nil {
					if err := cancelErrJob(ctx, jobTmp.ID(), dbTask); err != nil {
						logit.Context(ctx).WarnW("Cron.reloadTaskListTask.cancelErrJob.Err", err.Error())
					}
					continue
				}

				var jobNew gocron.Job
				jobNew, errNew = cron.Update(
					jobTmp.ID(),
					jobDefinitionNew,
					taskNew,
					getJobOptionList(ctx, jobTmp.ID(), dbTask)...,
				)
				if errNew != nil {
					if err := cancelErrJob(ctx, jobTmp.ID(), dbTask); err != nil {
						logit.Context(ctx).WarnW("Cron.reloadTaskListTask.cancelErrJob.Err", err.Error())
					}
					continue
				}
				logit.Context(ctx).DebugW("Cron.reloadTaskListTask", "Updated:"+jobNew.Name())
			}
		}

	})
	return
}

func cancelErrJob(ctx context.Context, taskUUID uuid.UUID, task db.GMTask) (err error) {
	if err = cron.RemoveJob(taskUUID); err != nil {
		return
	}
	task.Status = db.StatusPending
	task.UpdatedAt = uint64(time.Now().Unix())
	return global.DefaultDB.WithContext(ctx).Save(&task).Error
}
