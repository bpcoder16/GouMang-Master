package controller

import (
	"goumang-master/db"
	"goumang-master/errorcode"
	"goumang-master/global"

	"github.com/gin-gonic/gin"
)

type Task struct{}

func (t *Task) Config(ctx *gin.Context) {
	var dbNodeList []db.GMNode
	if err := global.DefaultDB.WithContext(ctx).Where(&db.GMNode{
		Status: db.StatusEnabled,
	}).Find(&dbNodeList).Error; err != nil {
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}
	returnSuccessJson(ctx, gin.H{
		"typeMethodList": []map[string]any{
			{"name": db.TaskMethodMap[db.NotChoose], "value": db.NotChoose},
		},
		"taskTypeList": []map[string]any{
			{"name": db.TaskTypeMap[db.NotChoose], "value": db.NotChoose},
			{"name": db.TaskTypeMap[db.TypeCron], "value": db.TypeCron},
			{"name": db.TaskTypeMap[db.TypeDuration], "value": db.TypeDuration},
			{"name": db.TaskTypeMap[db.TypeDurationRandom], "value": db.TypeDurationRandom},
			{"name": db.TaskTypeMap[db.TypeOneTimeJobStartDateTimes], "value": db.TypeOneTimeJobStartDateTimes},
		},
		"taskStatusList": []map[string]any{
			{"name": db.TaskStatusMap[db.NotChoose], "value": db.NotChoose},
			{"name": db.TaskStatusMap[db.StatusPending], "value": db.StatusPending},
			{"name": db.TaskStatusMap[db.StatusEnabled], "value": db.StatusEnabled},
			{"name": db.TaskStatusMap[db.StatusConfigError], "value": db.StatusConfigError},
			{"name": db.TaskStatusMap[db.StatusConfigExpired], "value": db.StatusConfigExpired},
			{"name": db.TaskStatusMap[db.StatusWorkerServiceError], "value": db.StatusWorkerServiceError},
		},
		"nodeList": func() []map[string]any {
			nodeList := make([]map[string]any, 0, len(dbNodeList)+1)
			nodeList = append(nodeList, map[string]any{
				"name": "全部", "value": db.NotChoose,
			})
			for _, node := range dbNodeList {
				nodeList = append(nodeList, map[string]any{
					"name":  node.Title,
					"value": node.ID,
				})
			}
			return nodeList
		}(),
	})
}

func (t *Task) List(ctx *gin.Context) {
	var req struct {
		Page     int    `json:"page" binding:"required,min=1"`
		PageSize int    `json:"page_size" binding:"required,min=1"`
		Title    string `json:"title"`
		Tag      string `json:"tag"`
		Method   int8   `json:"method"`
		NodeId   uint64 `json:"node_id"`
		Status   int8   `json:"status"`
	}

	if err := paramsValidator(ctx, &req); err != nil {
		return
	}

	where := make([]string, 0, 5)
	whereParams := make([]any, 0, 5)
	if len(req.Title) > 0 {
		where = append(where, "title LIKE ?")
		whereParams = append(whereParams, "%"+req.Title+"%")
	}
	if len(req.Tag) > 0 {
		where = append(where, "tag LIKE ?")
		whereParams = append(whereParams, "%"+req.Tag+"%")
	}
	if req.Method != 0 {
		where = append(where, "method = ?")
		whereParams = append(whereParams, req.Method)
	}
	if req.NodeId != 0 {
		where = append(where, "node_id = ?")
		whereParams = append(whereParams, req.NodeId)
	}

}
