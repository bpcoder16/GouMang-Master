package cron

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bpcoder16/Chestnut/v2/appconfig/env"
	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

var (
	scheduler gocron.Scheduler
	once      sync.Once
)

func lazyInit() {
	once.Do(func() {
		var err error
		// 创建一个新的调度器
		scheduler, err = gocron.NewScheduler(
			gocron.WithLocation(env.TimeLocation()),
		)
		if err != nil {
			panic("创建调度器失败:" + err.Error())
		}
	})
}

func Jobs() []gocron.Job {
	return scheduler.Jobs()
}

func NewJob(jobDefinition gocron.JobDefinition, task gocron.Task, jobOptions ...gocron.JobOption) (gocron.Job, error) {
	return scheduler.NewJob(jobDefinition, task, jobOptions...)
}

func RemoveJob(uuidStr uuid.UUID) error {
	return scheduler.RemoveJob(uuidStr)
}

func Update(uuidStr uuid.UUID, jobDefinition gocron.JobDefinition, task gocron.Task, jobOptions ...gocron.JobOption) (gocron.Job, error) {
	return scheduler.Update(uuidStr, jobDefinition, task, jobOptions...)
}

func Run(ctx context.Context) {
	lazyInit()

	scheduler.Start()
	logit.Context(ctx).InfoW("cron.Manager.Run", "Cron 调度器已启动")

	// 捕获系统信号以优雅地关闭调度器
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logit.Context(ctx).InfoW("cron.Manager.Run", "Context cancelled, Cron 调度器准备关闭...")
	case sig := <-sigChan:
		logit.Context(ctx).InfoW("cron.Manager.Run", fmt.Sprintf("Received signal: %v, Cron 调度器准备关闭...", sig))
	}

	logit.Context(ctx).InfoW("cron.Manager.Run", "Cron 调度器关闭中...")
	if err := scheduler.Shutdown(); err != nil {
		logit.Context(ctx).ErrorW("cron.Manager.Run", "Cron 调度器关闭失败:"+err.Error())
	}
	logit.Context(ctx).InfoW("cron.Manager.Run", "Cron 调度器关闭完成，已退出")
}
