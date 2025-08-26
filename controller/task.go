package controller

import (
	"goumang-master/db"
	"goumang-master/errorcode"
	"goumang-master/global"
	"sort"
	"strconv"
	"time"

	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/gin-gonic/gin"
)

type Task struct{}

// List 获取任务列表
func (t *Task) List(ctx *gin.Context) {
	// 获取查询参数
	title := ctx.Query("title")   // 任务名称
	tag := ctx.Query("tag")       // 任务标签
	method := ctx.Query("method") // 任务方式
	nodeID := ctx.Query("nodeID") // 任务节点
	status := ctx.Query("status") // 状态
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("pageSize", "10")

	// 参数验证
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 || pageSizeInt > 100 {
		pageSizeInt = 10
	}

	// 构建查询条件
	query := global.DefaultDB.WithContext(ctx).Model(&db.GMTask{}).Where("status != ?", db.StatusDeleted)

	// 添加条件查询
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if tag != "" {
		switch tag {
		case "NoTag":
			query = query.Where("tag = ?", "")
		default:
			query = query.Where("tag = ?", tag)
		}
	}

	if method != "" {
		if methodInt, err := strconv.Atoi(method); err == nil && methodInt != db.NotChoose {
			query = query.Where("method = ?", methodInt)
		}
	}

	if status != "" {
		if statusInt, err := strconv.Atoi(status); err == nil && statusInt != db.NotChoose {
			query = query.Where("status = ?", statusInt)
		}
	}

	// 节点查询，关联 gm_nodes_tasks 表
	if nodeID != "" {
		if nodeIDInt, err := strconv.ParseUint(nodeID, 10, 64); err == nil && nodeIDInt != db.NotChoose {
			query = query.Joins("JOIN gm_nodes_tasks ON gm_tasks.id = gm_nodes_tasks.task_id").Where("gm_nodes_tasks.node_id = ?", nodeIDInt)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logit.Context(ctx).ErrorW("getTaskList.count.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 分页查询
	var tasks []db.GMTask
	offset := (pageInt - 1) * pageSizeInt
	if err := query.Offset(offset).Limit(pageSizeInt).Order("created_at DESC").Find(&tasks).Error; err != nil {
		logit.Context(ctx).ErrorW("getTaskList.find.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 构建返回数据
	taskList := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		// 格式化下次运行时间
		var nextRunTimeStr string
		if task.NextRunTime > 0 {
			nextRunTimeStr = time.Unix(0, task.NextRunTime*int64(time.Millisecond)).Format(time.DateTime)
		} else {
			nextRunTimeStr = ""
		}

		// 格式化创建时间和更新时间
		createdAtStr := time.Unix(task.CreatedAt, 0).Format(time.DateTime)
		updatedAtStr := time.Unix(task.UpdatedAt, 0).Format(time.DateTime)

		taskList = append(taskList, map[string]any{
			"id":           task.ID,
			"uuid":         task.UUID,
			"title":        task.Title,
			"tag":          task.Tag,
			"desc":         task.Desc,
			"type":         task.Type,
			"typeName":     db.TaskTypeMap[task.Type],
			"expression":   task.Expression,
			"method":       task.Method,
			"methodName":   db.TaskMethodMap[task.Method],
			"methodParams": task.MethodParams,
			"nextRunTime":  nextRunTimeStr,
			"editable":     task.Editable,
			"status":       task.Status,
			"statusName":   db.TaskStatusMap[task.Status],
			"errorMessage": task.ErrorMessage,
			"createdAt":    createdAtStr,
			"updatedAt":    updatedAtStr,
		})
	}

	returnSuccessJson(ctx, gin.H{
		"total": total,
		"list":  taskList,
	})
}

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
		"name": "全部", "value": 0,
	})
	for _, node := range nodes {
		nodeList = append(nodeList, map[string]any{
			"name":  node.Title + " - " + node.IP + ":" + strconv.FormatUint(uint64(node.Port), 10),
			"value": node.ID,
		})
	}
	return nodeList
}
