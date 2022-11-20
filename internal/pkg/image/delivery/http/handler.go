package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/api/errors/convert"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ImageHandler struct {
	useCase image.UseCase
	logger  hardlogger.Logger
}

func NewImageHandler(usecase image.UseCase, logger hardlogger.Logger) *ImageHandler {
	return &ImageHandler{
		useCase: usecase,
		logger:  logger,
	}
}

func (h *ImageHandler) UploadImage(echoCtx echo.Context) error {
	ctx := h.logger.WithCaller(echoCtx.Request().Context())

	uid := echoCtx.QueryParam("uid")
	if uid == "" {
		err := errors.New("empty event id")
		h.logger.WithCtx(ctx).WithError(err).Errorf("failed to upload image for event %s", uid)
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	mf, err := echoCtx.MultipartForm()
	if err != nil {
		h.logger.WithCtx(ctx).WithError(err).Errorf("failed to get multipart form")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	err = h.useCase.UploadImage(ctx, mf.File, uid)
	if err != nil {
		h.logger.WithCtx(ctx).WithError(err).Errorf("failed to upload image for event %s", uid)
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.NoContent(http.StatusOK)
}

func (h *ImageHandler) UploadImageAlbum(echoCtx echo.Context) error {
	ctx := h.logger.WithCaller(echoCtx.Request().Context())

	form, err := echoCtx.MultipartForm()
	if err != nil {
		h.logger.WithCtx(ctx).WithError(err).Errorf("failed to get multipart form")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	respFromVk, err := h.useCase.UploadImageAlbum(ctx, form)
	if err != nil {
		h.logger.WithCtx(ctx).WithError(err).Errorf("error upload images")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, respFromVk)
}
