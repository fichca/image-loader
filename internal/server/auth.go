package server

import (
	"context"
	"encoding/json"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/response"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type authService interface {
	Authorize(ctx context.Context, login, password string) (string, error)
}

type authHandler struct {
	logger *logrus.Logger
	r      *chi.Mux
	as     authService
}

func NewAuthHandler(logger *logrus.Logger, as authService, r *chi.Mux) *authHandler {
	return &authHandler{
		logger: logger,
		r:      r,
		as:     as,
	}
}

func (ah *authHandler) RegisterAuthRoutes() {
	ah.r.Get("/user/auth", ah.HandleAuthorize)
}

// HandleAuthorize issues a JWT
//
//	@Summary      Authorize
//	@Description  Issue JWT
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param        user    body     dto.AuthUserDto  true  "authorize user"
//	@Success      200  {array}   response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      404  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /user/auth [get]
func (ah *authHandler) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	var user dto.AuthUserDto

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ah.handleError(err, http.StatusBadRequest, w)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ah.logger.Error(err)
		}
	}(r.Body)

	token, err := ah.as.Authorize(r.Context(), user.Login, user.Password)
	if err != nil {
		ah.handleError(err, http.StatusInternalServerError, w)
		return
	}

	b, err := response.ParseResponse(token, false)
	if err != nil {
		ah.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		ah.handleError(err, http.StatusInternalServerError, w)
		return
	}
}

func (ah *authHandler) handleError(err error, status int, w http.ResponseWriter) {
	ah.logger.Error(err)
	w.WriteHeader(status)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		ah.logger.Error(err)
	}
}
