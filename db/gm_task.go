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
	Type         int8   `gorm:"column:type"`          // 执行方式
	Expression   string `gorm:"column:expression"`    // 执行方式表达式
	Method       int8   `gorm:"column:method"`        // 任务方式
	MethodParams string `gorm:"column:method_params"` // 任务方式参数
	NextRunTime  int64  `gorm:"column:next_run_time"` // 下一次执行时间(时间戳:毫秒)
	Editable     int8   `gorm:"column:editable"`      // 是否可编辑 [当执行过一次后,就不可再编辑 Method&MethodParams]
	Status       int8   `gorm:"column:status"`        // 状态
	ErrorMessage string `gorm:"column:error_message"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime"` // 创建时间(时间戳:秒)
	UpdatedAt    int64  `gorm:"column:updated_at;autoUpdateTime"` // 更新时间(时间戳:秒)
}

func (GMTask) TableName() string {
	return "gm_tasks"
}

const (
	NotChoose = 0

	// 执行方式
	TypeOneTimeJobStartImmediately int8 = -1 // (系统) 立即执行
	TypeCron                       int8 = 1  // 定时任务 [* * * * * *]
	TypeDuration                   int8 = 2  // 时间间隔, 单位:毫秒 [1000]
	TypeDurationRandom             int8 = 3  // 随机时间间隔, 单位:毫秒 [2000, 5000]
	TypeOneTimeJobStartDateTimes   int8 = 4  // 指定时间 [2025-08-22 10:00:00,2025-09-22 22:00:00]

	// 任务方式
	MethodTest               int8 = -1 // (系统) 测试任务
	MethodReloadTaskList     int8 = -2 // (系统) 重新加载任务列表
	MethodInitJobNextRunTime int8 = -3 // (系统) 初始化所有任务下一次时间
	MethodShell              int8 = 0  // shell 命令

	// 是否可编辑状态
	EditableYes int8 = 1 // 可编辑
	EditableNo  int8 = 2 // 不可编辑

	// 运行状态
	RunStatusRunning int8 = 1 // 运行中
	RunStatusSuccess int8 = 2 // 已完成
	RunStatusFailure int8 = 3 // 失败

	// 状态
	StatusDeleted            int8 = -1 // 删除
	StatusPending            int8 = 1  // 待启用/下线
	StatusConfigError        int8 = 2  // 配置异常
	StatusConfigExpired      int8 = 3  // 配置超时下线
	StatusWorkerServiceError int8 = 4  // 下游服务异常
	StatusEnabled            int8 = 10 // 已启用
)

var (
	TaskTypeMap = map[int8]string{
		NotChoose:                      "全部",
		TypeOneTimeJobStartImmediately: "立即执行",
		TypeCron:                       "定时任务",
		TypeDuration:                   "时间间隔",
		TypeDurationRandom:             "随机时间间隔",
		TypeOneTimeJobStartDateTimes:   "指定时间",
	}
	TaskMethodMap = map[int8]string{
		NotChoose:                "全部",
		MethodTest:               "测试任务",
		MethodReloadTaskList:     "重新加载任务列表",
		MethodInitJobNextRunTime: "初始化所有任务下一次时间",
	}
	TaskStatusMap = map[int8]string{
		NotChoose:                "全部",
		StatusDeleted:            "删除",
		StatusPending:            "待启用/下线",
		StatusConfigError:        "配置异常",
		StatusConfigExpired:      "配置自动超时下线",
		StatusWorkerServiceError: "下游服务异常",
		StatusEnabled:            "已启用",
	}
)

func BuildSHA256(task GMTask) string {
	return utils.SHA265String(strings.Join([]string{
		strconv.FormatUint(uint64(task.Type), 10),
		task.Expression,
		strconv.FormatUint(uint64(task.Method), 10),
		task.MethodParams,
	}, "_"))
}
