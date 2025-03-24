package rest

import (
	"fmt"
	"net/http"
	"strings"

	"codepair-sinarmas/models"
	"codepair-sinarmas/pkg/logger"
	"codepair-sinarmas/pkg/serror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	// Ensure the serror package is correctly imported and available
)

func handleError(ctx *gin.Context, statusCode int, errx serror.SError) (result gin.H) {
	if statusCode == 0 || statusCode == http.StatusInternalServerError {
		logger.Err(errx)
		ctx.JSON(errx.Code(), models.ResponseError{
			Message: "Internal server error",
			Error:   errx.Error(),
		})
		return
	}

	ctx.JSON(errx.Code(), models.ResponseError{
		Message: errx.Error(),
	})
	return
}

func handleValidationError(ctx *gin.Context, validationErrors interface{}) (result gin.H) {
	ctx.JSON(http.StatusUnprocessableEntity, models.ResponseError{
		Message: "Validation error",
		Error:   validationErrors,
	})

	return
}

func BuildAndGetValidationMessage(err error) string {
	var validationMessages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			validationMessages = append(validationMessages, fmt.Sprintf("Field '%s' is required.", err.Field()))
		case "min":
			validationMessages = append(validationMessages, fmt.Sprintf("Field '%s' must be at least %s characters long.", err.Field(), err.Param()))
		case "max":
			validationMessages = append(validationMessages, fmt.Sprintf("Field '%s' must not exceed %s characters.", err.Field(), err.Param()))
		case "eqfield":
			validationMessages = append(validationMessages, fmt.Sprintf("Field '%s' must match '%s'.", err.Field(), err.Param()))
		default:
			validationMessages = append(validationMessages, fmt.Sprintf("Field '%s' failed validation on rule '%s'.", err.Field(), err.Tag()))
		}
	}

	return strings.Join(validationMessages, " ")
}
