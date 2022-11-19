package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/api/errors/convert"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/complaint"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ComplaintHandler struct {
	useCase complaint.Usecase
	logger  hardlogger.Logger
}

func NewComplaintHandler(usecase complaint.Usecase, logger hardlogger.Logger) *ComplaintHandler {
	return &ComplaintHandler{
		useCase: usecase,
		logger:  logger,
	}
}

func (ch *ComplaintHandler) Create(echoCtx echo.Context) error {
	ctx := ch.logger.WithCaller(echoCtx.Request().Context())

	var newComplaint models.Complaint
	if err := echoCtx.Bind(&newComplaint); err != nil {
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind new complaint data")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	if err := echoCtx.Validate(&newComplaint); err != nil {
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed validate complaint")
		return echoCtx.JSON(http.StatusBadRequest, errors.New("failed validate data").Error())
	}

	err := ch.useCase.Create(ctx, newComplaint)
	if err != nil {
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to create complaint")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).ErrorExplain())
	}

	return echoCtx.NoContent(http.StatusOK)
}
