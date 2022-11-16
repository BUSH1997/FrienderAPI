package validator

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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

func SourceEvent(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if val != models.SOURCE_EVENT_USER ||
		val != models.SOURSE_EVENT_GROUP ||
		val != models.SOURCE_EVENT_FORK_GROUP ||
		val != models.SOURCE_EVENT_VK {
		return false
	}
	return true
}
