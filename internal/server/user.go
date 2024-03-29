package server

import (
	"context"
	"encoding/json"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type userService interface {
	Add(ctx context.Context, user dto.UserDto) error
	GetById(ctx context.Context, id int) (dto.UserResponse, error)
	Update(ctx context.Context, user dto.UserDto) error
	DeleteById(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]dto.UserDto, error)
}

type userHandler struct {
	logger         *logrus.Logger
	r              *chi.Mux
	us             userService
	authMiddleware func(next http.Handler) http.Handler
}

func NewUserHandler(logger *logrus.Logger, us userService, r *chi.Mux, authMiddleware func(next http.Handler) http.Handler) *userHandler {
	return &userHandler{
		logger:         logger,
		r:              r,
		us:             us,
		authMiddleware: authMiddleware,
	}
}

func (uh *userHandler) RegisterUserRoutes() {
	uh.r.Post("/user/add", uh.HandleAddUser)

	uh.r.Group(func(r chi.Router) {
		r.Use(uh.authMiddleware)

		r.Get("/user/{userID}", uh.HandleGetByIdUser)
		r.Put("/user/update", uh.HandleUpdateUser)
		r.Delete("/user/{userID}", uh.HandleDeleteByIdUser)
		r.Get("/user", uh.HandleGetAllUsers)
	})
}

// HandleAddUser adds a new user
//
//	@Summary      AddUser
//	@Description  add a new user
//	@Tags         user
//	@Accept       json
//	@Produce      json
//	@Param        user    body     dto.UserDto  true  "add a new user"
//	@Success      200  {array}   response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      404  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /user/add [post]
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

// HandleGetByIdUser get user by id
//
//	@Summary        GetUserById
//	@Description    get user
//	@Tags            user
//	@Accept            json
//	@Produce        json
//	@Param            id    path        string    true    "get user by ID"
//	@Success        200    {array}        dto.UserResponse
//	@Failure        400    {object}    response.Response
//	@Failure        404    {object}    response.Response
//	@Failure        500    {object}    response.Response
//	@Router            /user/{userID} [get]
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

// HandleUpdateUser update user
//
//	@Summary        UpdateUser
//	@Description    update user
//	@Tags            user
//	@Accept            json
//	@Produce        json
//	@Param            user    body        dto.UserDto    true    "update user"
//	@Success        200
//	@Failure        400        {object}    response.Response
//	@Failure        404        {object}    response.Response
//	@Failure        500        {object}    response.Response
//	@Router            /user/update [put]
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

// HandleDeleteByIdUser delete a user
//
//	@Summary        DeleteUser
//	@Description    delete a user
//	@Tags            user
//	@Accept            json
//	@Produce        json
//	@Param            id    path        string    true    "delete user"
//	@Success        200
//	@Failure        400    {object}    response.Response
//	@Failure        404    {object}    response.Response
//	@Failure        500    {object}    response.Response
//	@Router            /user/{userID} [delete]
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

// HandleGetAllUsers get all users
//
//	@Summary        GetAllUser
//	@Description    get all users
//	@Tags           user
//	@Accept         json
//	@Produce        json
//	@Success        200    {array}        dto.UserDto
//	@Failure        400    {object}    response.Response
//	@Failure        404    {object}    response.Response
//	@Failure        500    {object}    response.Response
//	@Router            /user/ [get]
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
