package db

type GMNode struct {
	ID        uint64 `gorm:"column:id"`
	Title     string `gorm:"column:title"`                     // 节点标题
	IP        string `gorm:"column:ip"`                        // IP
	Port      uint   `gorm:"column:port"`                      // Port
	Remark    string `gorm:"column:remark"`                    // 备注
	Status    int8   `gorm:"column:status"`                    // 状态
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime"` // 创建时间(时间戳:秒)
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime"` // 更新时间(时间戳:秒)
}

func (GMNode) TableName() string {
	return "gm_nodes"
}
