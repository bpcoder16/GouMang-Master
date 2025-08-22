package tasks

import (
	"context"

	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/google/uuid"
)

// Job 执行前执行
func beforeJobRunsFunc(ctx context.Context) func(jobID uuid.UUID, jobName string) {
	return func(jobID uuid.UUID, jobName string) {
		//logit.Context(ctx).DebugW("BeforeJobRuns.jobID", jobID, "jobName", jobName)
	}
}

// Job 执行前执行，抛出 error 可打断执行
func beforeJobRunsSkipIfBeforeFuncErrorsFunc(ctx context.Context) func(jobID uuid.UUID, jobName string) error {
	return func(jobID uuid.UUID, jobName string) (err error) {
		//var task db.GMTask
		//if err = global.DefaultDB.WithContext(ctx).Where("uuid = ?", jobID.String()).First(&task).Error; err != nil {
		//	return
		//}
		//
		//var nodeList []db.GMNodesTasks
		//if err = global.DefaultDB.WithContext(ctx).Where("task_id = ?", task.ID).Find(&nodeList).Error; err != nil {
		//	return
		//}
		//
		//nowAt := time.Now().UnixNano() / 1e6
		//if len(nodeList) == 0 {
		//	err = global.DefaultDB.WithContext(ctx).Create(&db.GMTaskLog{
		//		TaskID:         task.ID,
		//		TaskTitle:      task.Title,
		//		TaskType:       task.Type,
		//		TaskExpression: task.Expression,
		//		NodeID:         0,
		//		StartedAt:      nowAt,
		//		RunStatus:      db.RunStatusRunning,
		//		CreatedAt:      nowAt,
		//		UpdatedAt:      nowAt,
		//	}).Error
		//} else {
		//	logList := make([]db.GMTaskLog, 0, len(nodeList))
		//	for _, node := range nodeList {
		//		logList = append(logList, db.GMTaskLog{
		//			TaskID:         task.ID,
		//			TaskTitle:      task.Title,
		//			TaskType:       task.Type,
		//			TaskExpression: task.Expression,
		//			NodeID:         node.ID,
		//			StartedAt:      nowAt,
		//			RunStatus:      db.RunStatusRunning,
		//			CreatedAt:      nowAt,
		//			UpdatedAt:      nowAt,
		//		})
		//	}
		//	err = global.DefaultDB.WithContext(ctx).Create(&logList).Error
		//}
		return
	}
}

// Job 执行后执行
func afterJobRunsFunc(ctx context.Context) func(jobID uuid.UUID, jobName string) {
	return func(jobID uuid.UUID, jobName string) {
		if job, err := GetJob(jobID.String()); err != nil {
			return
		} else {
			_ = updateDBTaskNextRunTime(ctx, job, true)
		}

		//
		//var logList []db.GMTaskLog
		//if err := global.DefaultDB.WithContext(ctx).
		//	Where("task_id = ? and run_status = ?", task.ID, db.RunStatusRunning).Find(&logList).Error; err != nil {
		//	return
		//}
		//
		//idList := make([]uint64, 0, len(logList))
		//for _, log := range logList {
		//	idList = append(idList, log.ID)
		//}
		//
		//nowAt := time.Now().UnixNano() / 1e6
		//
		//global.DefaultDB.WithContext(ctx).Where("id in ?", idList).Model(&db.GMTaskLog{}).Updates(map[string]interface{}{
		//	"run_status": db.RunStatusSuccess,
		//	"ended_at":   nowAt,
		//	"updated_at": nowAt,
		//})
		//
	}
}

// Job 执行后执行，接收 Error
func afterJobRunsWithErrorFunc(ctx context.Context) func(jobID uuid.UUID, jobName string, err error) {
	return func(jobID uuid.UUID, jobName string, err error) {
		logit.Context(ctx).ErrorW("AfterJobRunsWithError.jobID", jobID, "jobName", jobName, "error", err)
	}
}

// Job 执行后执行，接收 Panic
func afterJobRunsWithPanicFunc(ctx context.Context) func(jobID uuid.UUID, jobName string, recoverData any) {
	return func(jobID uuid.UUID, jobName string, recoverData any) {
		logit.Context(ctx).ErrorW("AfterJobRunsWithPanic.jobID", jobID, "jobName", jobName, "recoverData", recoverData)
	}
}
