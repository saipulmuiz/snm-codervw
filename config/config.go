package config

import (
	"os"

	"codepair-sinarmas/pkg/serror"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	DB     *gorm.DB
	Server *gin.Engine
}

func Init() (cfg Config) {
	Catch(cfg.InitTimezone())
	Catch(cfg.InitPostgres())
	Catch(cfg.InitService())

	return
}

func (cfg *Config) Start() (errx serror.SError) {
	cfg.Server.Run(os.Getenv("APP_PORT"))

	return
}

// func (cfg *Config) stop() {
// 	if cfg.APM != nil {
// 		cfg.APM.Shutdown()
// 	}
// }

func Catch(errx serror.SError) {
	if errx != nil {
		errx.Panic()
	}
}
