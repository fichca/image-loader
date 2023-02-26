package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fichca/image-loader/internal/dto"
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

type userHandler struct {
	listenURI string
	logger    *logrus.Logger
	r         chi.Router
	us        userService
}

func NewUserHandler(listenURI string, logger *logrus.Logger, us userService) *userHandler {
	return &userHandler{
		listenURI: listenURI,
		logger:    logger,
		r:         chi.NewRouter(),
		us:        us,
	}
}

func (uh *userHandler) RegisterRoutes() {
	uh.r.Get("/user/{userID}", uh.HandleGetByIdUser)
	uh.r.Post("/user/add", uh.HandleAddUser)
	uh.r.Put("/user/update", uh.HandleUpdateUser)
	uh.r.Delete("/user/delete/{userID}", uh.HandleDeleteByIdUser)
	uh.r.Get("/user", uh.HandleGetAllUsers)
}

func (uh *userHandler) StartServer() {
	srv := http.Server{
		Addr:    uh.listenURI,
		Handler: uh.r,
	}

	uh.logger.Info(fmt.Sprintf("server is running on port %v!", uh.listenURI))
	err := srv.ListenAndServe()
	if err != nil {
		uh.logger.Fatal(err)
	}
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
