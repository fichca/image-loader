package server

import (
	"context"
	"github.com/fichca/image-loader/internal/response"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type fileService interface {
	AddFile(ctx context.Context, filename string, file io.Reader) error
}

type fileHandler struct {
	logger *logrus.Logger
	r      *chi.Mux
	fs     fileService
}

func NewFileHandler(logger *logrus.Logger, fs fileService, r *chi.Mux) *fileHandler {
	return &fileHandler{
		logger: logger,
		r:      r,
		fs:     fs,
	}
}

func (fh *fileHandler) RegisterFileRoutes() {
	fh.r.Get("/image/add", fh.HandleAddFile)
}

func (fh *fileHandler) HandleAddFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("fileKey")
	if err != nil {
		fh.handleError(err, http.StatusBadRequest, w)
	}

	err = fh.fs.AddFile(r.Context(), header.Filename, file)
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
