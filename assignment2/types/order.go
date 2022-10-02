package types

import "time"

type OrderItem struct {
	ItemID      uint   `json:"lineItemId"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
}

type ReqWriteOrder struct {
	OrderedAt    time.Time   `json:"orderedAt"`
	CustomerName string      `json:"customerName"`
	Items        []OrderItem `json:"items"`
}
