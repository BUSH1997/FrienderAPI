package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type ImageHandler struct {
	useCase image.UseCase
}

func NewImageHandler(usecase image.UseCase) *ImageHandler {
	return &ImageHandler{
		useCase: usecase,
	}
}

func (h *ImageHandler) UploadImage(ctx echo.Context) error {
	uid := ctx.QueryParam("uid")
	if uid == "" {
		log.Error("Baduid")
		return ctx.NoContent(http.StatusBadRequest)
	}
	mf, err := ctx.MultipartForm()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err = h.useCase.UploadImage(ctx.Request().Context(), mf.File, uid)
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
