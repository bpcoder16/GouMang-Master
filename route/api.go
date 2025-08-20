package route

import (
	"goumang-master/db"
	"goumang-master/services/tasks"
	"net/http"

	"github.com/bpcoder16/Chestnut/v2/contrib/httphandler/gin"
	"github.com/bpcoder16/Chestnut/v2/core/utils"
	gin2 "github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Api() gin.Router {
	apiRouter := gin.NewDefaultRouter("/api")

	apiRouter.GET("/test", func(ctx *gin2.Context) {
		taskUUID := uuid.New()
		tasks.CreateJob(ctx, db.GMTask{
			UUID:   taskUUID.String(),
			Title:  "TestImmediately",
			Type:   db.TypeOneTimeJobStartImmediately,
			Method: db.MethodTest,
		})
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"uuid":   uuid.New(),
			"sha256": utils.SHA265String("1_16 */30 * * * *_1_"),
		})
	})
	//apiRouter.GET("/task-list", func(ctx *gin2.Context) {
	//	var taskList []db.GMTask
	//	err := sqlite.DefaultClient().WithContext(ctx).Find(&taskList).Error
	//
	//	ctx.JSON(http.StatusOK, gin2.H{
	//		"taskList": taskList,
	//		"err":      err,
	//	})
	//})
	//
	//apiRouter.GET("/add-job", func(ctx *gin2.Context) {
	//	taskUUID := uuid.New()
	//	_, err := cron.NewJob(
	//		gocron.DurationJob(time.Second*5),
	//		gocron.NewTask(func() {
	//			fmt.Println(time.Now().Format(time.RFC3339Nano) + " add-job")
	//		}),
	//		gocron.WithContext(ctx),
	//		gocron.WithName("add-job"),
	//		gocron.WithSingletonMode(gocron.LimitModeWait),
	//		gocron.WithTags("add-job", "test"),
	//		gocron.WithIdentifier(taskUUID),
	//	)
	//
	//	ctx.JSON(http.StatusOK, gin2.H{
	//		"uuid": taskUUID,
	//		"err":  err,
	//	})
	//})
	//
	//apiRouter.GET("/job-list", func(ctx *gin2.Context) {
	//	jobList := cron.Jobs()
	//	type jobShowStruct struct {
	//		ID      string
	//		LastRun time.Time
	//		Name    string
	//		NextRun time.Time
	//		Tags    []string
	//	}
	//	jobShowList := make([]jobShowStruct, 0, len(jobList))
	//	for _, job := range jobList {
	//		lastRun, _ := job.LastRun()
	//		nextRun, _ := job.NextRun()
	//		jobShowList = append(jobShowList, jobShowStruct{
	//			ID:      job.ID().String(),
	//			LastRun: lastRun,
	//			Name:    job.Name(),
	//			NextRun: nextRun,
	//			Tags:    job.Tags(),
	//		})
	//	}
	//
	//	ctx.JSON(http.StatusOK, gin2.H{
	//		"jobShowList": jobShowList,
	//	})
	//})

	return apiRouter
}
