package models

import "gorm.io/gorm"

type SocialMedia struct {
	gorm.Model
	Name   string `json:"name"`
	URL    string `json:"social_media_url"`
	UserID uint   `json:"user_id"`
}
