package models

import (
	"final_project/helpers"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;uniqueIndex" json:"username" validate:"required"`
	Email    string `gorm:"not null;uniqueIndex" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"password" validate:"required,min=6"`
	Age      uint   `gorm:"not null" json:"age" validate:"required,numeric,min=8"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	validate := validator.New()
	if err = validate.Struct(u); err != nil {
		return err
	}
	u.Password = helpers.HashPass(u.Password)
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	validate := validator.New()
	if err = validate.Struct(u); err != nil {
		return err
	}
	return nil
}
