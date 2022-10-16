package models

import "gorm.io/gorm"

type SocialMedia struct {
	gorm.Model
	Name   string `gorm:"not null" json:"name" validate:"required"`
	URL    string `gorm:"not null" json:"social_media_url" validate:"required"`
	UserID uint   `json:"user_id"`
	User   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
}
