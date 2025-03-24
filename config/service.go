package config

import (
	"codepair-sinarmas/pkg/serror"
	"codepair-sinarmas/service/handler/rest"
	"codepair-sinarmas/service/repository"
	"codepair-sinarmas/service/usecase"
)

func (cfg *Config) InitService() (errx serror.SError) {
	userRepo := repository.NewUserRepository(cfg.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)

	otpRepo := repository.NewOtpRepository(cfg.DB)
	otpUsecase := usecase.NewOtpUsecase(otpRepo, userRepo)

	route := rest.CreateHandler(
		userUsecase,
		otpUsecase,
	)

	cfg.Server = route

	return nil
}
