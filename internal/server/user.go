package server

import (
	"context"
	"encoding/json"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type userService interface {
	Add(ctx context.Context, user dto.UserDto) error
	GetById(ctx context.Context, id int) (dto.UserDto, error)
	Update(ctx context.Context, user dto.UserDto) error
	DeleteById(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]dto.UserDto, error)
}

type validateService interface {
	ValidateUser(ctx context.Context, user dto.AuthUserDto) error
}

type userHandler struct {
	logger     *logrus.Logger
	r          *chi.Mux
	us         userService
	vs         validateService
	jwtKeyword string
}

func NewUserHandler(logger *logrus.Logger, us userService, vs validateService, r *chi.Mux, jwtKeyword string) *userHandler {
	return &userHandler{
		logger:     logger,
		r:          r,
		us:         us,
		vs:         vs,
		jwtKeyword: jwtKeyword,
	}
}

func (uh *userHandler) RegisterUserRoutes() {
	uh.r.Post("/user/add", uh.HandleAddUser)

	uh.r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(uh.vs, uh.jwtKeyword, uh.logger))

		r.Get("/user/{userID}", uh.HandleGetByIdUser)
		r.Put("/user/update", uh.HandleUpdateUser)
		r.Delete("/user/delete/{userID}", uh.HandleDeleteByIdUser)
		r.Get("/user", uh.HandleGetAllUsers)
	})
}

func (uh *userHandler) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var user dto.UserDto

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uh.handleError(err, http.StatusBadRequest, w)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			uh.logger.Error(err)
		}
	}(r.Body)

	err = uh.us.Add(context.Background(), user)
	if err != nil {
		uh.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uh *userHandler) HandleGetByIdUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		uh.handleError(err, http.StatusBadRequest, w)
		return
	}

	user, err := uh.us.GetById(context.Background(), id)
	if err != nil {
		uh.handleError(err, http.StatusInternalServerError, w)
		return
	}

	b, err := json.Marshal(&user)
	if err != nil {
		uh.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		uh.logger.Error(err)
	}
}

func (uh *userHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.UserDto

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uh.handleError(err, http.StatusBadRequest, w)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			uh.logger.Error(err)
		}
	}(r.Body)

	err = uh.us.Update(r.Context(), user)
	if err != nil {
		uh.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (uh *userHandler) HandleDeleteByIdUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		uh.handleError(err, http.StatusBadRequest, w)
		return
	}

	err = uh.us.DeleteById(context.Background(), id)
	if err != nil {
		uh.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uh *userHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := uh.us.GetAll(context.Background())
	if err != nil {
		uh.handleError(err, http.StatusInternalServerError, w)
		return
	}

	b, err := json.Marshal(&users)
	if err != nil {
		uh.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		uh.logger.Error(err)
	}

}

func (uh *userHandler) handleError(err error, status int, w http.ResponseWriter) {
	uh.logger.Error(err)
	w.WriteHeader(status)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		uh.logger.Error(err)
	}
}
