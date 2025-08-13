package db

type GMTask struct {
	ID           uint64 `gorm:"column:id"`
	UUID         string `gorm:"column:uuid"`
	UserID       uint64 `gorm:"column:user_id"`       // 用户id
	Title        string `gorm:"column:title"`         // 任务标题
	Desc         string `gorm:"column:desc"`          // 任务描述
	Type         uint8  `gorm:"column:type"`          // 任务执行类型
	Expression   string `gorm:"column:expression"`    // 任务表达式
	Method       uint8  `gorm:"column:method"`        // 任务执行方式
	MethodParams string `gorm:"column:method_params"` // 任务执行方式参数
	Status       int8   `gorm:"column:status"`        // -1 删除 1 待启用 2 已启用
	CreatedAt    uint64 `gorm:"column:created_at"`    // 创建时间
	UpdatedAt    uint64 `gorm:"column:updated_at"`    // 更新时间
}

// TableName 自定义表名，GORM 默认是结构体名的复数形式
func (GMTask) TableName() string {
	return "gm_tasks"
}
