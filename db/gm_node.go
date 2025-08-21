package db

type GMNode struct {
	ID     uint64 `gorm:"column:id"`
	Title  string `gorm:"column:title"`  // 节点标题
	IP     string `gorm:"column:ip"`     // IP
	Port   int    `gorm:"column:port"`   // Port
	Remark string `gorm:"column:remark"` // 备注
	Status int8   `gorm:"column:status"` // -3 下游服务异常 1 待启用/下线 2 已启用
}

func (GMNode) TableName() string {
	return "gm_nodes"
}
