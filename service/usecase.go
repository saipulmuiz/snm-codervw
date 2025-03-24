package service

import (
	"codepair-sinarmas/models"
	"codepair-sinarmas/pkg/serror"
)

type UserUsecase interface {
	Register(request *models.RegisterUser) (user *models.User, errx serror.SError)
	Login(request *models.LoginUser) (res models.LoginResponse, errx serror.SError)
}

type OTPUsecase interface {
	RequestOtp(request *models.OTPRequest) (res *models.OTPResponse, errx serror.SError)
	ValidateOtp(request *models.OTPValidateRequest) (valid bool, errx serror.SError)
}
