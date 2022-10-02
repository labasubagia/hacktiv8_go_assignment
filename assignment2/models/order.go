package models

import (
	"time"
)

type Order struct {
	OrderID      uint      `gorm:"primaryKey" json:"order_id"`
	CustomerName string    `gorm:"not null;type:varchar(191)" json:"customer_name"`
	OrderedAt    time.Time `json:"ordered_at"`
	Items        []Item    `json:"items"`
}
