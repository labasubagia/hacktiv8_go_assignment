package models

type Item struct {
	ItemID      uint   `gorm:"primaryKey" json:"item_id"`
	ItemCode    string `gorm:"not null;type:varchar(191)" json:"item_code"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"order_id"`
}
