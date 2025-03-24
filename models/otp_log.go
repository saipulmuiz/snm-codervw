package models

import (
	"time"
)

type OTPLog struct {
	ID               int64  `gorm:"primaryKey"`
	UserID           int64  `gorm:"not null"`
	OTPCode          string `gorm:"not null"`
	NotificationType string `gorm:"not null"`
	Status           string `gorm:"not null"`
	CreatedAt        time.Time
	ExpiredAt        time.Time
}

type OTPRequest struct {
	UserID int64 `json:"user_id"`
}

type OTPValidateRequest struct {
	UserID int64  `json:"user_id"`
	OTP    string `json:"otp"`
}

type OTPResponse struct {
	UserID int64  `json:"user_id"`
	OTP    string `json:"otp"`
}
