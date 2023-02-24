package server

import (
	"encoding/json"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type userService interface {
	AddUser(user dto.UserDto) error
	GetUser(id int) (dto.UserDto, error)
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

func (s *userHandler) RegisterRoutes() {
	s.r.Get("/user/{userID}", s.HandleGetUser)
	s.r.Post("/user/add", s.HandleAddUser)
}

func (s *userHandler) StartServer() {
	srv := http.Server{
		Addr:    s.listenURI,
		Handler: s.r,
	}

	s.logger.Info("server is running!")
	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Fatal(err)
	}
}

func (s *userHandler) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var user dto.UserDto

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}(r.Body)

	err = s.us.AddUser(user)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *userHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	user, err := s.us.GetUser(id)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	b, err := json.Marshal(&user)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *userHandler) handleError(err error, status int, w http.ResponseWriter) {
	s.logger.Error(err)
	w.WriteHeader(status)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		s.logger.Error(err)
	}
}
