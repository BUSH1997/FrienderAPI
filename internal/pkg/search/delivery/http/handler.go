package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/search"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SearchHandler struct {
	useCase search.UseCase
	logger  *logrus.Logger
}

func NewSearchHandler(useCase search.UseCase, logger *logrus.Logger) *SearchHandler {
	return &SearchHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (sh *SearchHandler) Search(ctx echo.Context) error {
	var wordList models.WordList
	if err := ctx.Bind(&wordList); err != nil {
		sh.logger.WithError(err).Errorf("failed to bind wordlist")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	events, err := sh.useCase.Search(wordList.Words)
	if err != nil {
		sh.logger.WithError(err).Errorf("failed to get events by search")
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, events)
}
