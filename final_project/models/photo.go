package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
	UserID   uint   `json:"user_id"`
}
