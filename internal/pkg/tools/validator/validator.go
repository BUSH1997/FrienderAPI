package validator

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("Failed validate data"))
	}
	return nil
}

func TitleEvent(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if len(val) < 5 {
		return false
	}
	return true
}

func DescriptionEvent(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if len(val) < 5 {
		return false
	}
	return true
}
