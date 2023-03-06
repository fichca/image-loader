package middleware

import (
	"context"
	"errors"
	"github.com/fichca/image-loader/internal/constants"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type authService interface {
	ValidateUser(ctx context.Context, user dto.AuthUserDto) (int64, error)
}

func Auth(us authService, keyword string, logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user, err := initCredentials(r.Header, keyword)
			if err != nil {
				writeErr(err, logger, w)
				return
			}
			userId, err := us.ValidateUser(context.Background(), user)

			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.IdCtxKey, int(userId))
			if err != nil {
				writeErr(err, logger, w)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func initCredentials(header http.Header, keyword string) (dto.AuthUserDto, error) {
	tokenStr := header.Get("Authorization")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(keyword), nil
	})

	if err != nil {
		return dto.AuthUserDto{}, err
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		return dto.AuthUserDto{}, err
	}
	credentials := strings.Split(issuer, " ")
	if len(credentials) != 2 {
		return dto.AuthUserDto{}, errors.New("math: square root of negative number")
	}
	return dto.NewAuthUserDto(credentials[0], credentials[1]), nil
}

func writeErr(err error, l *logrus.Logger, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	l.Error(err)

	b, err := response.ParseResponse(err.Error(), true)
	if err != nil {
		l.Error(err)
	}

	_, err = w.Write(b)
	if err != nil {
		l.Error(err)
	}
}
