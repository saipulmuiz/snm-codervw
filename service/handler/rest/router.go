package rest

import (
	"os"

	api "codepair-sinarmas/service"
	middlewares "codepair-sinarmas/service/middleware"

	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	userUsecase api.UserUsecase
	otpUsecase  api.OTPUsecase
}

func CreateHandler(
	userUsecase api.UserUsecase,
	otpUsecase api.OTPUsecase,
) *gin.Engine {
	obj := Handler{
		userUsecase: userUsecase,
		otpUsecase:  otpUsecase,
	}

	var maxSize int64 = 1024 * 1024 * 10 //10 MB
	logger := log.New()
	r := gin.Default()
	mainRouter := r.Group("/v1")

	gin.SetMode(gin.DebugMode)
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	corsconfig := cors.DefaultConfig()
	corsconfig.AllowAllOrigins = true
	corsconfig.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsconfig))
	r.Use(limits.RequestSizeLimiter(maxSize))
	r.Use(middlewares.ErrorHandler(logger))

	mainRouter.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	mainRouter.POST("/register", obj.Register)
	mainRouter.POST("/login", obj.login)
	mainRouter.POST("/otp-request", obj.RequestOTP)
	mainRouter.POST("/otp-validate", obj.ValidateOTP)

	authorizedRouter := mainRouter.Group("/")
	authorizedRouter.Use(middlewares.Auth())
	{
	}

	return r
}
