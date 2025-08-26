package controller

import (
	"goumang-master/db"
	"goumang-master/errorcode"
	"goumang-master/global"
	"sort"

	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/gin-gonic/gin"
)

type Task struct{}

// Config 获取任务配置相关信息
func (t *Task) Config(ctx *gin.Context) {
	// 获取标签列表
	tags, err := t.getTagList(ctx)
	if err != nil {
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 获取节点列表
	nodes, err := t.getNodeList(ctx)
	if err != nil {
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 构建标签列表
	tagList := buildTagList(tags)

	// 构建任务方式列表
	methodList := buildMethodList()

	// 构建任务状态列表
	statusList := buildStatusList()

	// 构建节点列表
	nodeList := buildNodeList(nodes)

	returnSuccessJson(ctx, gin.H{
		"tagList":    tagList,
		"methodList": methodList,
		"statusList": statusList,
		"nodeList":   nodeList,
	})
}



// getTagList 获取标签列表
func (t *Task) getTagList(ctx *gin.Context) (tagList []string, err error) {
	if err = global.DefaultDB.WithContext(ctx).Model(&db.GMTask{}).
		Where("status != ?", db.StatusDeleted).
		Distinct("tag").
		Pluck("tag", &tagList).Error; err != nil {
		logit.Context(ctx).ErrorW("getTagList.error", err)
		return
	}

	// 按字母顺序排序
	sort.Strings(tagList)
	return
}

// getNodeList 获取节点列表
func (t *Task) getNodeList(ctx *gin.Context) (nodeList []db.GMNode, err error) {
	err = global.DefaultDB.WithContext(ctx).Model(&db.GMNode{}).Where("status != ?", db.StatusDeleted).Find(&nodeList).Error
	if err != nil {
		logit.Context(ctx).ErrorW("getNodeList.error", err)
		return
	}
	return
}

// buildTagList 构建标签列表
func buildTagList(tags []string) []map[string]any {
	tagList := make([]map[string]any, 0, len(tags)+1)
	tagList = append(tagList, map[string]any{
		"name": "全部", "value": "",
	}, map[string]any{
		"name": "无标签", "value": "NoTag",
	})
	for _, tag := range tags {
		tagList = append(tagList, map[string]any{
			"name":  tag,
			"value": tag,
		})
	}
	return tagList
}

// buildMethodList 构建任务方式列表
func buildMethodList() []map[string]any {
	// 按固定顺序定义任务方式
	methodOrder := []int8{
		db.NotChoose,
		db.MethodTest,
		db.MethodReloadTaskList,
		db.MethodInitJobNextRunTime,
	}

	methodList := make([]map[string]any, 0, len(methodOrder))

	// 按固定顺序添加任务方式
	for _, value := range methodOrder {
		if name, exists := db.TaskMethodMap[value]; exists {
			methodList = append(methodList, map[string]any{
				"name":  name,
				"value": value,
			})
		}
	}

	return methodList
}

// buildStatusList 构建任务状态列表
func buildStatusList() []map[string]any {
	// 按固定顺序定义任务状态
	statusOrder := []int8{
		db.NotChoose,
		db.StatusPending,
		db.StatusConfigError,
		db.StatusConfigExpired,
		db.StatusWorkerServiceError,
		db.StatusEnabled,
	}

	statusList := make([]map[string]any, 0, len(statusOrder))

	// 按固定顺序添加任务状态
	for _, value := range statusOrder {
		if name, exists := db.TaskStatusMap[value]; exists {
			statusList = append(statusList, map[string]any{
				"name":  name,
				"value": value,
			})
		}
	}

	return statusList
}

// buildNodeList 构建节点列表
func buildNodeList(nodes []db.GMNode) []map[string]any {
	nodeList := make([]map[string]any, 0, len(nodes)+1)
	nodeList = append(nodeList, map[string]any{
		"name": "全部", "value": "",
	})
	for _, node := range nodes {
		nodeList = append(nodeList, map[string]any{
			"name":  node.Title,
			"value": node.ID,
		})
	}
	return nodeList
}
