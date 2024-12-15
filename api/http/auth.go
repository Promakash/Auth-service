package http

import (
	"auth_service/api/http/types"
	"auth_service/domain"
	pkghttp "auth_service/pkg/http"
	"auth_service/usecases"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
)

type AuthHandler struct {
	logger  slog.Logger
	service usecases.User
}

func NewAuthHandler(logger slog.Logger, service usecases.User) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: service,
	}
}

func (h *AuthHandler) WithAuthHandlers() pkghttp.RouterOption {
	return func(r chi.Router) {
		r.Route(
			"/tokens",
			func(r chi.Router) {
				pkghttp.AddHandler(r.Get, "/generate", h.generateHandler)
				pkghttp.AddHandler(r.Post, "/refresh", h.refreshHandler)
			},
		)
	}
}

// generateHandler обрабатывает запрос на генерацию токенов.
//
// @Summary Генерация токенов
// @Description Возвращает пару Access и Refresh токенов для указанного пользователя.
// @Tags Tokens
// @Produce json
// @Param uuid query string true "UUID пользователя"
// @Success 200 {object} types.GenerateResponse
// @Failure 404 {object} http.ErrorResponse "Пользователя с данным UUID нет в системе"
// @Failure 500 {object} http.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tokens/generate [get]
func (h *AuthHandler) generateHandler(r *http.Request) pkghttp.Response {
	req, err := types.CreateGenerateRequest(r)
	if err != nil {
		return pkghttp.BadRequest(err)
	}
	IP := pkghttp.ReadUserIP(r)
	aToken, rToken, err := h.service.GenTokens(req.UserID, IP)
	return domain.HandleError(err, types.GenerateResponse{AccessToken: aToken, RefreshToken: rToken})
}

// refreshHandler обрабатывает запрос на обновление токенов.
//
// @Summary Обновление токенов
// @Description Обновляет пару Access и Refresh токенов на основе переданного Refresh токена.
// @Tags Tokens
// @Accept json
// @Produce json
// @Param body body types.RefreshRequest true "Тело запроса"
// @Success 200 {object} types.RefreshResponse
// @Failure 401 {object} http.ErrorResponse "Недействительный/просроченный токен"
// @Failure 500 {object} http.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tokens/refresh [post]
func (h *AuthHandler) refreshHandler(r *http.Request) pkghttp.Response {
	req, err := types.CreateRefreshRequest(r)
	if err != nil {
		return pkghttp.BadRequest(err)
	}
	IP := pkghttp.ReadUserIP(r)
	aToken, rToken, err := h.service.RefreshTokens(req.RefreshToken, IP)
	return domain.HandleError(err, types.RefreshResponse{AccessToken: aToken, RefreshToken: rToken})
}
