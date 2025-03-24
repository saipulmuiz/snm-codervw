package service

import (
	"codepair-sinarmas/models"
)

type UserRepository interface {
	Register(user *models.User) (*models.User, error)
	GetUserByID(userId int64) (user *models.User, err error)
	GetUserByEmail(email string) (user *models.User, err error)
}

type OtpRepository interface {
	SaveOTP(otpLog *models.OTPLog) (*models.OTPLog, error)
	GetOtpByUserID(userID int64) (otp *models.OTPLog, err error)
	GetOtpByUserIDAndCode(userID int64, otpCode string) (otp *models.OTPLog, err error)
}
