package controller

import (
	"errors"
	"goumang-master/db"
	"goumang-master/errorcode"
	"goumang-master/global"
	"strconv"
	"strings"
	"time"

	"github.com/bpcoder16/Chestnut/v2/core/utils"
	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Node struct{}

// List 获取节点列表
func (n *Node) List(ctx *gin.Context) {
	// 获取查询参数
	title := strings.TrimSpace(ctx.Query("title")) // 节点名称
	ip := strings.TrimSpace(ctx.Query("ip"))       // 主机IP
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
	query := global.DefaultDB.WithContext(ctx).Model(&db.GMNode{}).Where("status = ?", db.StatusEnabled)

	// 添加条件查询
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if ip != "" {
		query = query.Where("ip LIKE ?", "%"+ip+"%")
	}

	// 获取总数
	var total int64
	if errT := query.Count(&total).Error; errT != nil {
		logit.Context(ctx).ErrorW("getNodeList.count.error", errT)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 分页查询
	var dbNodes []db.GMNode
	offset := (pageInt - 1) * pageSizeInt
	if err = query.Offset(offset).Limit(pageSizeInt).Order("created_at ASC").Find(&dbNodes).Error; err != nil {
		logit.Context(ctx).ErrorW("getNodeList.find.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 构建返回数据
	nodeList := make([]map[string]any, 0, len(dbNodes))
	for _, node := range dbNodes {
		// 格式化创建时间和更新时间
		createdAtStr := time.Unix(node.CreatedAt, 0).Format(time.DateTime)
		updatedAtStr := time.Unix(node.UpdatedAt, 0).Format(time.DateTime)

		nodeList = append(nodeList, map[string]any{
			"id":        node.ID,
			"title":     node.Title,
			"ip":        node.IP,
			"port":      node.Port,
			"remark":    node.Remark,
			"createdAt": createdAtStr,
			"updatedAt": updatedAtStr,
		})
	}

	returnSuccessJson(ctx, gin.H{
		"total": total,
		"list":  nodeList,
	})
}

// CreateNodeRequest 创建节点请求参数
type CreateNodeRequest struct {
	Title  string `json:"title" binding:"required"` // 节点标题
	IP     string `json:"ip" binding:"required"`    // IP地址
	Port   uint   `json:"port" binding:"required"`  // 端口
	Remark string `json:"remark"`                   // 备注
}

type EditNodeRequest struct {
	ID                uint64 `json:"id" binding:"required"`
	CreateNodeRequest `json:",inline"`
}

func (n *Node) createUpdateCheck(ctx *gin.Context, req *CreateNodeRequest, excludeId uint64) (err error) {
	// 参数处理
	req.Title = strings.TrimSpace(req.Title)
	req.IP = strings.TrimSpace(req.IP)
	req.Remark = strings.TrimSpace(req.Remark)

	// 参数验证
	if req.Title == "" {
		err = errors.New("节点标题不能为空")
		return
	}

	if !utils.CheckIPv4Valid(req.IP) {
		err = errors.New("IP地址格式有问题")
		return
	}

	if !utils.CheckPortValid(req.Port) {
		err = errors.New("端口格式有问题")
		return
	}

	// 检查节点标题是否已存在
	if errDB := global.DefaultDB.WithContext(ctx).
		Where("title = ? and status = ? and id != ?", req.Title, db.StatusEnabled, excludeId).
		First(&db.GMNode{}).Error; !errors.Is(errDB, gorm.ErrRecordNotFound) {
		err = errors.New("节点标题已存在，不可重复")
		return
	}

	// 检查 IP&PORT 是否已存在
	if errDB := global.DefaultDB.WithContext(ctx).
		Where("ip = ? and port = ? and status = ? and id != ?", req.IP, req.Port, db.StatusEnabled, excludeId).
		First(&db.GMNode{}).Error; !errors.Is(errDB, gorm.ErrRecordNotFound) {
		err = errors.New("IP&PORT 已存在，不可重复")
		return
	}

	return
}

// Create 创建节点
func (n *Node) Create(ctx *gin.Context) {
	var req CreateNodeRequest
	if err := paramsValidator(ctx, &req); err != nil {
		return
	}

	if err := n.createUpdateCheck(ctx, &req, 0); err != nil {
		returnErrJson(ctx, errorcode.ErrParams, err.Error())
		return
	}

	// 创建节点
	newNode := db.GMNode{
		Title:  req.Title,
		IP:     req.IP,
		Port:   req.Port,
		Remark: req.Remark,
		Status: db.StatusEnabled,
	}

	if err := global.DefaultDB.WithContext(ctx).Create(&newNode).Error; err != nil {
		logit.Context(ctx).ErrorW("createNode.create.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	returnSuccessJson(ctx, gin.H{
		"id": newNode.ID,
	})
}

// Edit 编辑节点
func (n *Node) Edit(ctx *gin.Context) {
	var req EditNodeRequest
	if err := paramsValidator(ctx, &req); err != nil {
		return
	}

	// 检查节点是否存在
	node, err := n.getNodeByID(ctx, req.ID)
	if err != nil {
		returnErrJson(ctx, errorcode.ErrParams, err.Error())
		return
	}

	createNodeRequest := req.CreateNodeRequest
	if errC := n.createUpdateCheck(ctx, &createNodeRequest, req.ID); errC != nil {
		returnErrJson(ctx, errorcode.ErrParams, errC.Error())
		return
	}

	// 更新节点信息
	node.Title = createNodeRequest.Title
	node.IP = createNodeRequest.IP
	node.Port = createNodeRequest.Port
	node.Remark = createNodeRequest.Remark

	if err = global.DefaultDB.WithContext(ctx).Save(&node).Error; err != nil {
		logit.Context(ctx).ErrorW("editNode.save.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	returnSuccessJson(ctx, gin.H{})
}

func (n *Node) getNodeByID(ctx *gin.Context, id uint64) (node db.GMNode, err error) {
	err = global.DefaultDB.WithContext(ctx).Where("id = ? and status = ?", id, db.StatusEnabled).
		First(&node).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("节点不存在")
			return
		}
		logit.Context(ctx).ErrorW("getNodeByID.error", err)
	}
	return
}

// Detail 获取节点详情
func (n *Node) Detail(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		returnErrJson(ctx, errorcode.ErrParams, "无效的节点ID")
		return
	}

	var node db.GMNode
	if node, err = n.getNodeByID(ctx, id); err != nil {
		returnErrJson(ctx, errorcode.ErrParams, err.Error())
		return
	}

	// 格式化时间
	createdAtStr := time.Unix(node.CreatedAt, 0).Format("2006-01-02 15:04:05")
	updatedAtStr := time.Unix(node.UpdatedAt, 0).Format("2006-01-02 15:04:05")

	returnSuccessJson(ctx, gin.H{
		"id":        node.ID,
		"title":     node.Title,
		"ip":        node.IP,
		"port":      node.Port,
		"remark":    node.Remark,
		"createdAt": createdAtStr,
		"updatedAt": updatedAtStr,
	})
}

// Delete 删除节点
func (n *Node) Delete(ctx *gin.Context) {
	var req struct {
		ID uint64 `json:"id" validate:"required"`
	}
	if err := paramsValidator(ctx, &req); err != nil {
		return
	}

	// 检查节点是否存在
	if _, err := n.getNodeByID(ctx, req.ID); err != nil {
		returnErrJson(ctx, errorcode.ErrParams, err.Error())
		return
	}

	// 检查是否有任务使用该节点，且该节点是任务的唯一节点
	var taskIDs []uint64
	if err := global.DefaultDB.WithContext(ctx).Model(&db.GMNodesTasks{}).Where("node_id = ?", req.ID).Pluck("task_id", &taskIDs).Error; err != nil {
		logit.Context(ctx).ErrorW("deleteNode.getTaskIDs.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 检查每个任务是否只有这一个节点
	for _, taskID := range taskIDs {
		var nodeCount int64
		if err := global.DefaultDB.WithContext(ctx).Model(&db.GMNodesTasks{}).Where("task_id = ?", taskID).Count(&nodeCount).Error; err != nil {
			logit.Context(ctx).ErrorW("deleteNode.countNodes.error", err)
			returnErrJson(ctx, errorcode.ErrServiceException)
			return
		}

		if nodeCount == 1 {
			// 获取任务信息
			var task db.GMTask
			if err := global.DefaultDB.WithContext(ctx).Where("id = ?", taskID).First(&task).Error; err == nil {
				returnErrJson(ctx, errorcode.ErrParams, "无法删除节点，任务'"+task.Title+"'仅绑定了该节点")
				return
			}
			returnErrJson(ctx, errorcode.ErrParams, "无法删除节点，存在任务仅绑定了该节点")
			return
		}
	}

	// 在事务中删除节点和关联的任务关系
	if err := global.DefaultDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除节点
		if err := tx.Model(&db.GMNode{}).Where("id = ?", req.ID).
			Update("status", db.StatusDeleted).Error; err != nil {
			return err
		}
		return tx.Where("node_id = ?", req.ID).Delete(&db.GMNodesTasks{}).Error
	}); err != nil {
		logit.Context(ctx).ErrorW("deleteNode.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	returnSuccessJson(ctx, gin.H{})
}

// GetNodeTasks 获取使用指定节点的任务列表
func (n *Node) GetNodeTasks(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		returnErrJson(ctx, errorcode.ErrParams, "无效的节点ID")
		return
	}

	// 检查节点是否存在
	var node db.GMNode
	if node, err = n.getNodeByID(ctx, id); err != nil {
		returnErrJson(ctx, errorcode.ErrParams, err.Error())
		return
	}

	// 获取使用该节点的任务列表
	var tasks []db.GMTask
	if err = global.DefaultDB.WithContext(ctx).
		Joins("JOIN gm_nodes_tasks ON gm_tasks.id = gm_nodes_tasks.task_id").
		Where("gm_nodes_tasks.node_id = ? AND gm_tasks.status != ?", id, db.StatusDeleted).
		Find(&tasks).Error; err != nil {
		logit.Context(ctx).ErrorW("getNodeTasks.findTasks.error", err)
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 构建返回数据
	taskList := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		taskList = append(taskList, buildShowTask(task, []uint64{}))
	}

	returnSuccessJson(ctx, gin.H{
		"nodeInfo": gin.H{
			"id":     node.ID,
			"title":  node.Title,
			"ip":     node.IP,
			"port":   node.Port,
			"remark": node.Remark,
		},
		"tasks":     taskList,
		"taskTotal": len(taskList),
	})
}
