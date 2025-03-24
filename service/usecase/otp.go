package usecase

import (
	"net/http"
	"time"

	"codepair-sinarmas/models"
	"codepair-sinarmas/pkg/serror"
	"codepair-sinarmas/service"
	api "codepair-sinarmas/service"
	"codepair-sinarmas/service/helper"

	"gorm.io/gorm"
)

type OtpUsecase struct {
	otpRepo  api.OtpRepository
	userRepo api.UserRepository
}

func NewOtpUsecase(
	otpRepo api.OtpRepository,
	userRepo api.UserRepository,
) service.OTPUsecase {
	return &OtpUsecase{
		otpRepo:  otpRepo,
		userRepo: userRepo,
	}
}

func (u *OtpUsecase) RequestOtp(request *models.OTPRequest) (res *models.OTPResponse, errx serror.SError) {
	user, err := u.userRepo.GetUserByID(request.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errx = serror.Newi(http.StatusNotFound, "User not found")
			return
		}

		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][RequestOtp] Failed to get user by email, [userID: %s]", request.UserID)
		return
	}

	otp, err := u.otpRepo.GetOtpByUserID(user.UserID)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][RequestOtp] Failed to get otp by userID, [userID: %s]", request.UserID)
		return
	}

	if otp.Status == "created" {
		res = &models.OTPResponse{
			UserID: user.UserID,
			OTP:    otp.OTPCode,
		}

		return
	}

	otpArgs := &models.OTPLog{
		UserID:           user.UserID,
		NotificationType: models.NotificationTypeEmail,
		OTPCode:          helper.GenerateOtpCode(),
		Status:           "created",
		CreatedAt:        time.Now(),
		ExpiredAt:        time.Now().Add(time.Minute * 2),
	}

	otp, err = u.otpRepo.SaveOTP(otpArgs)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][RequestOtp] Failed to save otp, [userID: %s]", request.UserID)
		return
	}

	res = &models.OTPResponse{
		UserID: user.UserID,
		OTP:    otp.OTPCode,
	}

	return
}

func (u *OtpUsecase) ValidateOtp(request *models.OTPValidateRequest) (valid bool, errx serror.SError) {
	_, err := u.userRepo.GetUserByID(request.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errx = serror.Newi(http.StatusNotFound, "User not found")
			return
		}

		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][ValidateOtp] Failed to get otp by userID, [userID: %s]", request.UserID)
		return
	}

	otpDB, err := u.otpRepo.GetOtpByUserIDAndCode(request.UserID, request.OTP)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errx = serror.Newi(http.StatusNotFound, "Invalid OTP code")
			return
		}
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][ValidateOtp] Failed to get otp by userID, [userID: %s]", request.UserID)
		return
	}

	if time.Now().After(otpDB.ExpiredAt) {
		errx = serror.Newi(http.StatusBadRequest, "OTP has expired")
		return
	}

	if otpDB.Status == "validated" {
		errx = serror.Newi(http.StatusBadRequest, "OTP has been validated")
		return
	}

	err = u.otpRepo.UpdateStatusOtpByUserIDAndCode(otpDB.ID, "validated")
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][ValidateOtp] Failed to update otp status, [userID: %s]", request.UserID)
		return
	}

	valid = true
	return
}
