package route

import (
	"fmt"
	"goumang-master/db"
	"net/http"
	"time"

	"github.com/bpcoder16/Chestnut/v2/contrib/cron"
	"github.com/bpcoder16/Chestnut/v2/contrib/httphandler/gin"
	"github.com/bpcoder16/Chestnut/v2/sqlite"
	gin2 "github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func Api() gin.Router {
	apiRouter := gin.NewDefaultRouter("/api")

	apiRouter.GET("/task-list", func(ctx *gin2.Context) {
		var taskList []db.GMTask
		err := sqlite.DB().WithContext(ctx).Find(&taskList).Error

		ctx.JSON(http.StatusOK, gin2.H{
			"taskList": taskList,
			"err":      err,
		})
	})

	apiRouter.GET("/add-job", func(ctx *gin2.Context) {
		taskUUID := uuid.New()
		_, err := cron.NewJob(
			gocron.DurationJob(time.Second*5),
			gocron.NewTask(func() {
				fmt.Println(time.Now().Format(time.RFC3339Nano) + " add-job")
			}),
			gocron.WithContext(ctx),
			gocron.WithName("add-job"),
			gocron.WithSingletonMode(gocron.LimitModeWait),
			gocron.WithTags("add-job", "test"),
			gocron.WithIdentifier(taskUUID),
		)

		ctx.JSON(http.StatusOK, gin2.H{
			"uuid": taskUUID,
			"err":  err,
		})
	})

	apiRouter.GET("/job-list", func(ctx *gin2.Context) {
		jobList := cron.Jobs()
		type jobShowStruct struct {
			ID      string
			LastRun time.Time
			Name    string
			NextRun time.Time
			Tags    []string
		}
		jobShowList := make([]jobShowStruct, 0, len(jobList))
		for _, job := range jobList {
			lastRun, _ := job.LastRun()
			nextRun, _ := job.NextRun()
			jobShowList = append(jobShowList, jobShowStruct{
				ID:      job.ID().String(),
				LastRun: lastRun,
				Name:    job.Name(),
				NextRun: nextRun,
				Tags:    job.Tags(),
			})
		}

		ctx.JSON(http.StatusOK, gin2.H{
			"jobShowList": jobShowList,
		})
	})

	return apiRouter
}
