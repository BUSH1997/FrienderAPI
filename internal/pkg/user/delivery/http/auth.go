package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"github.com/BUSH1997/FrienderAPI/internal/api/errors/convert"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (h *UserHandler) Auth(echoCtx echo.Context) error {
	ctx := h.Logger.WithCaller(echoCtx.Request().Context())

	var authParams models.AuthParams
	if err := echoCtx.Bind(&authParams); err != nil {
		h.Logger.WithCtx(ctx).WithError(err).Errorf("failed to bind auth params data")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	rowURL := url.Values{
		"vk_app_id":  []string{strconv.Itoa(int(authParams.AppID))},
		"vk_user_id": []string{strconv.Itoa(int(authParams.UserID))},
		"vk_time":    []string{strconv.Itoa(int(authParams.Time))},
	}

	encodedURL := rowURL.Encode()
	hashAlg := hmac.New(sha256.New, []byte(h.AuthSecret))
	hashAlg.Write([]byte(encodedURL))
	hash := hashAlg.Sum(nil)

	if string(hash) != authParams.Sign {
		err := errors.New("invalid sign")
		h.Logger.WithCtx(ctx).WithError(err).Errorf("failed to check sign")
		return echoCtx.JSON(http.StatusUnauthorized, convert.DeliveryError(err).Error())
	}

	newAuthToken, err := h.UserUseCase.GenerateAuthToken(echoCtx.Request().Context())
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, errors.Wrap(err, "failed to get token"))
	}

	authCookie := &http.Cookie{
		Name:     "session_id",
		Value:    newAuthToken.Value,
		HttpOnly: true,
		Expires:  time.Unix(newAuthToken.Expires, 0),
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Path:     "/",
	}

	echoCtx.SetCookie(authCookie)

	userAgent := echoCtx.Request().Header.Get("User-Agent")
	userIP := echoCtx.Request().Header.Get("X-Forwarded-For")

	newRefreshToken, err := h.UserUseCase.UpdateRefresh(echoCtx.Request().Context(), models.FingerPrintData{
		UserAgent: userAgent,
		UserIP:    userIP,
	})
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to update refresh token"))
	}

	cookie := &http.Cookie{
		Name:     "refresh",
		Value:    newRefreshToken.Value,
		HttpOnly: true,
		Expires:  time.Unix(newRefreshToken.Expires, 0),
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Path:     "/",
	}

	echoCtx.SetCookie(cookie)

	return echoCtx.JSON(http.StatusOK, "authorized successfully")
}
