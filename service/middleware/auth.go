package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"codepair-sinarmas/models"
	"codepair-sinarmas/service/helper"

	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "Token Not Found",
				Error:   "Unauthorized",
			})

			return
		}

		bearer := strings.HasPrefix(token, "Bearer")

		if !bearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "Bearer Not Found",
				Error:   "Unauthorized",
			})

			return
		}
		tokenStr := strings.Split(token, "Bearer ")[1]

		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "Token STR Not Found",
				Error:   "Unauthorized",
			})

			return
		}

		claims, err := helper.VerifyToken(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: err.Error(),
				Error:   "Unauthorized",
			})

			log.Errorln("ERROR:", err)
			return
		}

		var data = claims.(jwt.MapClaims)
		userId := data["id"].(float64)
		strUserId := strconv.FormatFloat(userId, 'f', -1, 64)

		ctx.Set("user_id", strUserId)
		ctx.Set("name", data["name"])
		ctx.Set("email", data["email"])
		ctx.Set("exp", data["exp"])
		ctx.Set("exp_date", data["exp_date"])

		if data["exp_date"] == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "Invalid token",
				Error:   "Unauthorized",
			})

			return
		}

		timeNow := time.Now()
		expiredTime := data["exp_date"].(string)

		parsed, err := time.Parse(time.RFC3339, expiredTime)
		if err != nil {
			log.Errorln("ERROR:", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: err.Error(),
				Error:   "Unauthorized",
			})

			return
		}

		if timeNow.After(parsed) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "Token has expired, please login again",
				Error:   "Unauthorized",
			})

			return
		}

		ctx.Next()
	}
}
