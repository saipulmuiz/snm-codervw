package usecase

import (
	"net/http"

	"codepair-sinarmas/models"
	"codepair-sinarmas/pkg/serror"
	api "codepair-sinarmas/service"
	"codepair-sinarmas/service/helper"

	"gorm.io/gorm"
)

type UserUsecase struct {
	userRepo api.UserRepository
}

func NewUserUsecase(
	userRepo api.UserRepository,
) api.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) Register(request *models.RegisterUser) (user *models.User, errx serror.SError) {
	userArgs := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	userCheck, err := u.userRepo.GetUserByEmail(request.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][Register] Failed to get user by email, [email: %s]", request.Email)
		return
	}

	if userCheck.UserID != 0 {
		errx = serror.Newi(http.StatusBadRequest, "Email already registered")
		return
	}

	user, err = u.userRepo.Register(userArgs)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][Register] Failed to register user, [email: %s]", request.Email)
		return
	}

	return
}

func (u *UserUsecase) Login(request *models.LoginUser) (res models.LoginResponse, errx serror.SError) {
	userDB, err := u.userRepo.GetUserByEmail(request.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errx = serror.Newi(http.StatusNotFound, "User not found")
			return
		}

		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][Login] Failed to get user by email, [email: %s]", request.Email)
		return
	}

	accountMatch := helper.ComparePassword([]byte(userDB.Password), []byte(request.Password))
	if !accountMatch {
		errx = serror.Newi(http.StatusBadRequest, "Password does not match")
		return
	}

	token := helper.GenerateToken(userDB.UserID, userDB.Email, userDB.Name)

	res = models.LoginResponse{
		Token: token,
		User:  *userDB,
	}

	return
}
