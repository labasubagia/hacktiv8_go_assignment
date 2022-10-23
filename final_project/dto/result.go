package dto

import (
	"time"
)

type User struct {
	ID        *uint      `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Email     *string    `json:"email,omitempty"`
	Username  *string    `json:"username,omitempty"`
	Password  *string    `json:"password,omitempty"`
	Age       *uint      `json:"age,omitempty"`
}

type Photo struct {
	ID        *uint      `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Title     *string    `json:"title,omitempty"`
	URL       *string    `json:"photo_url,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	User      *User      `json:"user,omitempty"`
	Caption   *string    `json:"caption,omitempty"`
}

type Comment struct {
	ID        *uint      `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	PhotoID   *uint      `json:"photo_id,omitempty"`
	Message   *string    `json:"message,omitempty"`
	User      *User      `json:"user,omitempty"`
	Photo     *Photo     `json:"photo,omitempty"`
}

type SocialMedia struct {
	ID        *uint      `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Name      *string    `json:"name,omitempty"`
	URL       *string    `json:"social_media_url,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	User      *User      `json:"user,omitempty"`
}

type FieldValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Response struct {
	Data            any                    `json:"data"`
	ValidationError []FieldValidationError `json:"validation_error"`
	Error           string                 `json:"error"`
	Message         string                 `json:"message"`
}
