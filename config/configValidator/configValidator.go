package configValidator

import (
	customValidator "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/validator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ConfigValidator(router *echo.Echo) {
	val := validator.New()
	val.RegisterValidation("custom_title", customValidator.TitleEvent)
	val.RegisterValidation("custom_description", customValidator.DescriptionEvent)
	router.Validator = &customValidator.CustomValidator{Validator: val}
}
