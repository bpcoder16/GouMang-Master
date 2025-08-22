package db

type GMTaskLog struct {
	ID             uint64 `gorm:"column:id"`
	TaskID         uint64 `gorm:"column:task_id"`
	TaskTitle      string `gorm:"column:task_title"`      // 任务标题
	TaskType       uint8  `gorm:"column:task_type"`       // 任务执行类型	1:cron(* * * * * *) 2:duration[单位为毫秒](10) 3:durationRandom[单位为毫秒](10,50)
	TaskExpression string `gorm:"column:task_expression"` // 任务表达式
	NodeID         uint64 `gorm:"column:node_id"`
	StartedAt      int64  `gorm:"column:started_at"`                // 开始时间(时间戳:毫秒)
	EndedAt        int64  `gorm:"column:ended_at"`                  // 结束时间(时间戳:毫秒)
	RunStatus      int8   `gorm:"column:run_status"`                // 1 进行中 2 成功 2 失败
	CreatedAt      int64  `gorm:"column:created_at;autoCreateTime"` // 创建时间(时间戳:秒)
	UpdatedAt      int64  `gorm:"column:updated_at;autoUpdateTime"` // 更新时间(时间戳:秒)
}

func (GMTaskLog) TableName() string {
	return "gm_task_logs"
}
