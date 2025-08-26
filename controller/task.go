package controller

import (
	"errors"
	"goumang-master/db"
	"goumang-master/errorcode"
	"goumang-master/global"
	"goumang-master/services/tasks"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bpcoder16/Chestnut/v2/core/utils"
	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct{}

// List 获取任务列表
func (t *Task) List(ctx *gin.Context) {
	// 获取查询参数
	title := strings.TrimSpace(ctx.Query("title")) // 任务名称
	tag := strings.TrimSpace(ctx.Query("tag"))     // 任务标签
	method := ctx.Query("method")                  // 任务方式
	nodeID := ctx.Query("nodeID")                  // 任务节点
	status := ctx.Query("status")                  // 状态
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("pageSize", "10")

	// 参数验证
	var pageInt, pageSizeInt int
	var err error
	if pageInt, err = strconv.Atoi(page); err != nil || pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt, err = strconv.Atoi(pageSize); err != nil || pageSizeInt < 1 || pageSizeInt > 100 {
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
		if methodInt, errM := strconv.Atoi(method); errM == nil && methodInt != db.NotChoose {
			query = query.Where("method = ?", methodInt)
		}
	}

	if status != "" {
		if statusInt, errS := strconv.Atoi(status); errS == nil && statusInt != db.NotChoose {
			query = query.Where("status = ?", statusInt)
		}
	}

	// 节点查询，关联 gm_nodes_tasks 表
	if nodeID != "" {
		if nodeIDInt, errN := strconv.ParseUint(nodeID, 10, 64); errN == nil && nodeIDInt != db.NotChoose {
			query = query.Joins("JOIN gm_nodes_tasks ON gm_tasks.id = gm_nodes_tasks.task_id").Where("gm_nodes_tasks.node_id = ?", nodeIDInt)
		}
	}

	// 获取总数
	var total int64
	if errT := query.Count(&total).Error; errT != nil {
		logit.Context(ctx).ErrorW("getTaskList.count.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 分页查询
	var dbTasks []db.GMTask
	offset := (pageInt - 1) * pageSizeInt
	if err = query.Offset(offset).Limit(pageSizeInt).Order("created_at DESC").Find(&dbTasks).Error; err != nil {
		logit.Context(ctx).ErrorW("getTaskList.find.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 构建返回数据
	taskList := make([]map[string]any, 0, len(dbTasks))
	for _, task := range dbTasks {
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

// CreateTaskRequest 创建任务请求参数
type CreateTaskRequest struct {
	Title        string   `json:"title" binding:"required"`      // 任务标题
	Tag          string   `json:"tag"`                           // 标签
	Desc         string   `json:"desc"`                          // 任务描述
	Type         int8     `json:"type" binding:"required"`       // 执行方式
	Expression   string   `json:"expression" binding:"required"` // 执行方式表达式
	Method       int8     `json:"method" binding:"required"`     // 任务方式
	MethodParams string   `json:"methodParams"`                  // 任务方式参数
	NodeIDs      []uint64 `json:"nodeIDs"`                       // 节点ID列表
}

// Create 创建任务
func (t *Task) Create(ctx *gin.Context) {
	var req CreateTaskRequest
	if err := paramsValidator(ctx, &req); err != nil {
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	req.Tag = strings.TrimSpace(req.Tag)
	req.Desc = strings.TrimSpace(req.Desc)
	req.Expression = strings.TrimSpace(req.Expression)
	req.MethodParams = strings.TrimSpace(req.MethodParams)
	req.NodeIDs = utils.RemoveDuplicates(req.NodeIDs)

	if req.Tag == "NoTag" || req.Tag == "System" {
		returnErrJson(ctx, errorcode.ErrParams, "标签不能设置为 NoTag Or System")
		return
	}

	if len(req.Title) == 0 {
		returnErrJson(ctx, errorcode.ErrParams, "任务标题不能为空")
		return
	}

	if err := global.DefaultDB.WithContext(ctx).Where("title = ?", req.Title).First(&db.GMTask{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		returnErrJson(ctx, errorcode.ErrParams, "任务标题已存在，不可重复")
		return
	}

	// 验证执行方式
	if _, exists := db.TaskTypeMap[req.Type]; !exists || req.Type == db.NotChoose {
		returnErrJson(ctx, errorcode.ErrParams, "无效的执行方式")
		return
	}

	// 验证任务方式
	if _, exists := db.TaskMethodMap[req.Method]; !exists || req.Method == db.NotChoose {
		returnErrJson(ctx, errorcode.ErrParams, "无效的任务方式")
		return
	}

	var reqExpressionErr error
	switch req.Type {
	case db.TypeCron:
		reqExpressionErr = tasks.IsValidCrontabExpression(ctx, req.Expression)
	case db.TypeDuration:
		_, reqExpressionErr = tasks.IsValidDurationExpression(ctx, req.Expression)
	case db.TypeDurationRandom:
		_, _, reqExpressionErr = tasks.IsValidDurationRandomExpression(ctx, req.Expression)
	case db.TypeOneTimeJobStartDateTimes:
		_, reqExpressionErr = tasks.IsValidOneTimeJobStartDateTimesExpression(ctx, req.Expression)
	}
	if reqExpressionErr != nil {
		returnErrJson(ctx, errorcode.ErrParams, "无效的执行方式表达式")
		return
	}

	// 创建任务对象
	task := db.GMTask{
		UUID:         uuid.New().String(),
		UserID:       1, // 后续增加 UserID
		Title:        req.Title,
		Tag:          req.Tag,
		Desc:         req.Desc,
		Type:         req.Type,
		Expression:   req.Expression,
		Method:       req.Method,
		MethodParams: req.MethodParams,
		NextRunTime:  0, // 初始为0，后续计算
		Editable:     db.EditableYes,
		Status:       db.StatusPending,
		ErrorMessage: "",
	}

	// 计算SHA256
	task.SHA256 = db.BuildSHA256(task)

	// 验证节点ID有效性
	if len(req.NodeIDs) > 0 {
		var validNodeCount int64
		if err := global.DefaultDB.WithContext(ctx).Model(&db.GMNode{}).
			Where("id IN ?", req.NodeIDs).
			Count(&validNodeCount).Error; err != nil {
			logit.Context(ctx).ErrorW("validateNodeIDs.error", err)
			returnErrJson(ctx, errorcode.ErrServiceException)
			return
		}
		if int(validNodeCount) != len(req.NodeIDs) {
			returnErrJson(ctx, errorcode.ErrParams, "存在无效的节点ID")
			return
		}
	}

	if errT := global.DefaultDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&task).Error; err != nil {
			return err
		}
		nodeTaskList := make([]db.GMNodesTasks, 0, len(req.NodeIDs))
		for _, nodeID := range req.NodeIDs {
			nodeTaskList = append(nodeTaskList, db.GMNodesTasks{
				NodeID: nodeID,
				TaskID: task.ID,
			})
		}
		if len(nodeTaskList) > 0 {
			return tx.Create(&nodeTaskList).Error
		}
		return nil
	}); errT != nil {
		logit.Context(ctx).ErrorW("createTask.error", errT)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 返回创建的任务信息
	returnSuccessJson(ctx, gin.H{
		"id":   task.ID,
		"uuid": task.UUID,
	})
}

func (t *Task) getTaskByID(ctx *gin.Context, id uint64) (task db.GMTask, err error) {
	err = global.DefaultDB.WithContext(ctx).Where("id = ? AND status != ?", id, db.StatusDeleted).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			returnErrJson(ctx, errorcode.ErrParams, "任务不存在")
			return
		}
		logit.Context(ctx).ErrorW("getTaskByID.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	} else {
		if task.UserID == 0 {
			err = errors.New("当前任务不可操作")
			returnErrJson(ctx, errorcode.ErrParams, "当前任务不可操作")
			return
		}
	}
	return
}

func (t *Task) buildShowTask(task db.GMTask, nodeIDs []uint64) map[string]any {
	// 格式化时间
	nextRunTimeStr := ""
	if task.NextRunTime > 0 {
		nextRunTimeStr = time.Unix(task.NextRunTime/1000, 0).Format("2006-01-02 15:04:05")
	}

	createdAtStr := ""
	if task.CreatedAt > 0 {
		createdAtStr = time.Unix(task.CreatedAt, 0).Format("2006-01-02 15:04:05")
	}

	updatedAtStr := ""
	if task.UpdatedAt > 0 {
		updatedAtStr = time.Unix(task.UpdatedAt, 0).Format("2006-01-02 15:04:05")
	}

	return map[string]any{
		"id":           task.ID,
		"uuid":         task.UUID,
		"sha256":       task.SHA256,
		"userID":       task.UserID,
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
		"nodeIDs":      nodeIDs,
	}
}

// Detail 获取任务详情
func (t *Task) Detail(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		returnErrJson(ctx, errorcode.ErrParams, "任务ID不能为空")
		return
	}

	taskID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		returnErrJson(ctx, errorcode.ErrParams, "无效的任务ID")
		return
	}

	var task db.GMTask
	task, err = t.getTaskByID(ctx, taskID)
	if err != nil {
		return
	}

	// 获取关联的节点列表
	var nodesTasks []db.GMNodesTasks
	if err = global.DefaultDB.WithContext(ctx).Where("task_id = ?", taskID).Find(&nodesTasks).Error; err != nil {
		logit.Context(ctx).ErrorW("getTaskNodes.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	nodeIDs := make([]uint64, 0, len(nodesTasks))
	for _, nt := range nodesTasks {
		nodeIDs = append(nodeIDs, nt.NodeID)
	}

	taskShow := t.buildShowTask(task, nodeIDs)

	// 返回任务详情
	returnSuccessJson(ctx, taskShow)
}

func (t *Task) ImmediatelyRun(ctx *gin.Context) {
	var req struct {
		ID uint64 `json:"id" validate:"required"`
	}
	if err := paramsValidator(ctx, &req); err != nil {
		return
	}

	// 检查任务是否存在且未删除
	task, err := t.getTaskByID(ctx, req.ID)
	if err != nil {
		return
	}

	task.Type = db.TypeOneTimeJobStartImmediately
	_, err = tasks.CreateJob(ctx, global.DefaultDB.WithContext(ctx), task)
	if err != nil {
		logit.Context(ctx).ErrorW("createJob.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}
	returnSuccessJson(ctx, gin.H{})
}

// Delete 删除任务
func (t *Task) Delete(ctx *gin.Context) {
	var req struct {
		ID uint64 `json:"id" validate:"required"`
	}
	if err := paramsValidator(ctx, &req); err != nil {
		return
	}

	// 检查任务是否存在且未删除
	task, err := t.getTaskByID(ctx, req.ID)
	if err != nil {
		return
	}

	// 在事务中删除任务和关联的节点关系
	if errT := global.DefaultDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 软删除任务（设置状态为删除）
		if errT := tx.Model(&db.GMTask{}).Where("id = ?", task.ID).Updates(map[string]any{
			"status":        db.StatusDeleted,
			"next_run_time": 0,
		}).Error; errT != nil {
			return errT
		}
		_ = tasks.RemoveJobForDBTask(ctx, task)
		return nil
	}); errT != nil {
		logit.Context(ctx).ErrorW("deleteTask.error", errT)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	returnSuccessJson(ctx, gin.H{})
}

// Enable 启用任务
func (t *Task) Enable(ctx *gin.Context) {
	var req struct {
		ID uint64 `json:"id" validate:"required"`
	}
	if err := paramsValidator(ctx, &req); err != nil {
		return
	}

	// 检查任务是否存在且未删除
	task, err := t.getTaskByID(ctx, req.ID)
	if err != nil {
		return
	}

	// 检查任务是否已经是启用状态
	if task.Status == db.StatusEnabled {
		returnErrJson(ctx, errorcode.ErrParams, "任务已经是启用状态")
		return
	}

	// 更新任务状态为启用
	if errT := global.DefaultDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 软删除任务（设置状态为删除）
		if errT := tx.Model(&db.GMTask{}).Where("id = ?", req.ID).Update("status", db.StatusEnabled).Error; errT != nil {
			return errT
		}
		_, errC := tasks.CreateJob(ctx, tx, task)
		return errC
	}); errT != nil {
		logit.Context(ctx).ErrorW("enableTask.error", errT)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	returnSuccessJson(ctx, gin.H{})
}

// Disable 下线任务
func (t *Task) Disable(ctx *gin.Context) {
	var req struct {
		ID uint64 `json:"id" validate:"required"`
	}
	if err := paramsValidator(ctx, &req); err != nil {
		return
	}

	// 检查任务是否存在且未删除
	task, err := t.getTaskByID(ctx, req.ID)
	if err != nil {
		return
	}

	// 检查任务是否已经是待启用/下线状态
	if task.Status != db.StatusEnabled {
		returnErrJson(ctx, errorcode.ErrParams, "任务当前不是启用状态")
		return
	}

	// 更新任务状态为待启用/下线
	if errT := global.DefaultDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 软删除任务（设置状态为删除）
		if errT := tx.Model(&db.GMTask{}).Where("id = ?", req.ID).Updates(map[string]any{
			"status":        db.StatusPending,
			"next_run_time": 0,
		}).Error; errT != nil {
			return errT
		}
		_ = tasks.RemoveJobForDBTask(ctx, task)
		return nil
	}); errT != nil {
		logit.Context(ctx).ErrorW("disableTask.error", errT)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	returnSuccessJson(ctx, gin.H{})
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

	// 构建执行方式列表
	typeList := buildTypeList()

	// 构建任务方式列表
	methodList := buildMethodList()

	// 构建任务状态列表
	statusList := buildStatusList()

	// 构建节点列表
	nodeList := buildNodeList(nodes)

	returnSuccessJson(ctx, gin.H{
		"tagList":    tagList,
		"typeList":   typeList,
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

// buildTypeList 构建执行方式列表
func buildTypeList() []map[string]any {
	typeList := make([]map[string]any, 0, len(db.TaskTypeMap))

	// 按照预定义顺序添加执行方式
	typeOrder := []int8{
		db.TypeCron,
		db.TypeDuration,
		db.TypeDurationRandom,
		db.TypeOneTimeJobStartDateTimes,
	}

	for _, typeValue := range typeOrder {
		if typeName, exists := db.TaskTypeMap[typeValue]; exists {
			typeList = append(typeList, map[string]any{
				"name":  typeName,
				"value": typeValue,
			})
		}
	}

	return typeList
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
