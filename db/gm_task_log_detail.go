package db

type GMTaskLogDetail struct {
	ID        uint64 `gorm:"column:id"`
	TaskLogID uint64 `gorm:"column:task_log_id"`
	Content   string `gorm:"column:content"`
	CreatedAt int64  `gorm:"column:created_at"` // 创建时间
}

func (GMTaskLogDetail) TableName() string {
	return "gm_task_log_detail"
}
