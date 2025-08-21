package db

type GMNodesTasks struct {
	ID     uint64 `gorm:"column:id"`
	NodeID uint64 `gorm:"column:node_id"`
	TaskID uint64 `gorm:"column:task_id"`
}

func (GMNodesTasks) TableName() string {
	return "gm_nodes_tasks"
}
