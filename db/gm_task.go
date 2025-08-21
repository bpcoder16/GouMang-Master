package db

import (
	"strconv"
	"strings"

	"github.com/bpcoder16/Chestnut/v2/core/utils"
)

type GMTask struct {
	ID           uint64 `gorm:"column:id"`
	UUID         string `gorm:"column:uuid"`
	SHA256       string `gorm:"column:sha256"`
	UserID       uint64 `gorm:"column:user_id"`       // 用户id, 0 为 系统用户
	Title        string `gorm:"column:title"`         // 任务标题
	Tag          string `gorm:"column:tag"`           // 标签
	Desc         string `gorm:"column:desc"`          // 任务描述
	Type         uint8  `gorm:"column:type"`          // 任务执行类型	1:cron(* * * * * *) 2:duration[单位为毫秒](10) 3:durationRandom[单位为毫秒](10,50)
	Expression   string `gorm:"column:expression"`    // 任务表达式
	Method       uint8  `gorm:"column:method"`        // 任务执行方式 1 自有任务(刷新用户配置的任务列表)
	MethodParams string `gorm:"column:method_params"` // 任务执行方式参数
	NextRunTime  int64  `gorm:"column:next_run_time"` // 下一次执行时间(毫秒)
	Editable     int8   `gorm:"column:editable"`      // 是否可编辑 [当执行过一次后,就不可再编辑 Method&MethodParams]
	Status       int8   `gorm:"column:status"`        // -4 删除 -3 下游服务异常 -2 配置超时下线 -1 配置异常 1 待启用/下线 2 已启用
	ErrorMessage string `gorm:"column:error_message"`
	CreatedAt    int64  `gorm:"column:created_at"` // 创建时间(秒)
	UpdatedAt    int64  `gorm:"column:updated_at"` // 更新时间(秒)
}

func (GMTask) TableName() string {
	return "gm_tasks"
}

const (
	EditableYes int8 = 1
	EditableNo  int8 = 2

	RunStatusRunning int8 = 1
	RunStatusSuccess int8 = 2
	RunStatusFailure int8 = 3

	StatusDeleted            int8 = -4 // 删除
	StatusWorkerServiceError int8 = -3 // 下游服务异常
	StatusConfigExpired      int8 = -2 // 配置超时下线
	StatusConfigError        int8 = -1 // 配置异常
	StatusPending            int8 = 1  // 待启用/下线
	StatusEnabled            int8 = 2  // 已启用

	TypeCron                       uint8 = 1
	TypeDuration                   uint8 = 2
	TypeDurationRandom             uint8 = 3
	TypeOneTimeJobStartImmediately uint8 = 4
	TypeOneTimeJobStartDateTimes   uint8 = 5

	MethodTest           uint8 = 0
	MethodReloadTaskList uint8 = 1
	MethodShell          uint8 = 2
)

func BuildSHA256(task GMTask) string {
	return utils.SHA265String(strings.Join([]string{
		strconv.FormatUint(uint64(task.Type), 10),
		task.Expression,
		strconv.FormatUint(uint64(task.Method), 10),
		task.MethodParams,
	}, "_"))
}
