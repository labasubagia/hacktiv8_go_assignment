package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title" validate:"required"`
	Caption string `json:"caption"`
	URL     string `gorm:"not null" json:"photo_url" validate:"required"`
	UserID  uint   `json:"user_id"`
	User    *User  `json:"user"`
}

func (p *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	validate := validator.New()
	if err = validate.Struct(p); err != nil {
		return err
	}
	return nil
}

func (p *Photo) BeforeUpdate(tx *gorm.DB) (err error) {
	validate := validator.New()
	if err = validate.Struct(p); err != nil {
		return err
	}
	return nil
}
