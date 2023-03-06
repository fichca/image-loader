package server

import (
	"context"
	"fmt"
	"github.com/fichca/image-loader/internal/constants"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/response"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
)

type fileService interface {
	AddImage(ctx context.Context, image dto.Image) error
}

type fileHandler struct {
	logger         *logrus.Logger
	r              *chi.Mux
	fs             fileService
	authMiddleware func(next http.Handler) http.Handler
}

func NewFileHandler(logger *logrus.Logger, fs fileService, r *chi.Mux, authMiddleware func(next http.Handler) http.Handler) *fileHandler {
	return &fileHandler{
		logger:         logger,
		r:              r,
		fs:             fs,
		authMiddleware: authMiddleware,
	}
}

func (fh *fileHandler) RegisterFileRoutes() {

	fh.r.Group(func(r chi.Router) {
		r.Use(fh.authMiddleware)
		r.Post("/image/add", fh.HandleAddFile)
	})
}

func (fh *fileHandler) HandleAddFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("fileKey")
	if err != nil {
		fh.handleError(err, http.StatusBadRequest, w)
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fh.handleError(err, http.StatusBadRequest, w)
		}
	}(file)
	userID, err := userIDFromCtx(r.Context())
	if err != nil {
		fh.handleError(err, http.StatusInternalServerError, w)
		return
	}

	err = fh.fs.AddImage(r.Context(), dto.Image{
		UserID:    userID,
		Name:      header.Filename,
		Data:      file,
		Extension: ".jpg",
	})

	if err != nil {
		fh.handleError(err, http.StatusInternalServerError, w)
	}

	w.WriteHeader(http.StatusOK)
}

func (fh *fileHandler) handleError(err error, status int, w http.ResponseWriter) {
	fh.logger.Error(err)
	w.WriteHeader(status)

	b, err := response.ParseResponse(err.Error(), true)
	if err != nil {
		fh.logger.Error(err)
	}

	_, err = w.Write(b)
	if err != nil {
		fh.logger.Error(err)
	}
}

func userIDFromCtx(ctx context.Context) (int, error) {
	idAny := ctx.Value(constants.IdCtxKey)

	id, ok := idAny.(int)
	if !ok {
		return 0, fmt.Errorf("couldn't cast user id from context")
	}

	return id, nil
}
