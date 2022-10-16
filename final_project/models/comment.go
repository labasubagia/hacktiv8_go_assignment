package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	PhotoID uint   `json:"photo_id"`
	Message string `gorm:"not null" json:"message" validate:"required"`
	User    *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Photo   *Photo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"photo"`
}

func (u *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	validate := validator.New()
	if err = validate.Struct(u); err != nil {
		return err
	}
	return nil
}

func (u *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	validate := validator.New()
	if err = validate.Struct(u); err != nil {
		return err
	}
	return nil
}
