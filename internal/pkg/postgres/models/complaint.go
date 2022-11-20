package models

type Complaint struct {
	ID          uint   `gorm:"column:id"`
	Initiator   int64  `gorm:"column:initiator"`
	Item        string `gorm:"column:item"`
	ItemUID     string `gorm:"column:item_uid"`
	TimeCreated int64  `gorm:"column:time_created"`
	IsProcessed bool   `gorm:"column:is_processed"`
	Reason      string `gorm:"column:reason"`
}

func (Complaint) TableName() string {
	return "complaints"
}
