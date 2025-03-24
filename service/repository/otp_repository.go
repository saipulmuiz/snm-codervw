package repository

import (
	"codepair-sinarmas/models"
	api "codepair-sinarmas/service"

	"gorm.io/gorm"
)

type otpRepo struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) api.OtpRepository {
	return &otpRepo{
		db: db,
	}
}

func (u *otpRepo) SaveOTP(otpLog *models.OTPLog) (*models.OTPLog, error) {
	return otpLog, u.db.Create(&otpLog).Error
}

func (u *otpRepo) GetOtpByUserID(userID int64) (otp *models.OTPLog, err error) {
	err = u.db.Where("user_id = ?", userID).First(&otp).Error
	if err == gorm.ErrRecordNotFound {
		return &models.OTPLog{}, nil
	}

	return otp, err
}

func (u *otpRepo) GetOtpByUserIDAndCode(userID int64, otpCode string) (otp *models.OTPLog, err error) {
	return otp, u.db.Where("user_id = ? AND  otp_code = ?", userID, otpCode).First(&otp).Error
}
