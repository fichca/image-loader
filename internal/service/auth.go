package service

import (
	"context"
	"fmt"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type authRepository interface {
	CheckAuth(ctx context.Context, login, password string) (entity.User, error)
}

type authService struct {
	repo       authRepository
	jwtKeyword string
}

func NewAuthService(repo authRepository, jwtKeyword string) *authService {
	return &authService{
		repo:       repo,
		jwtKeyword: jwtKeyword,
	}
}

func (a authService) Authorize(ctx context.Context, login, password string) (string, error) {
	user, err := a.repo.CheckAuth(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to authorize user: %w", err)
	}
	now := time.Now()
	issuer := fmt.Sprintf("%s %s", login, password)
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   "authorized",
		Audience:  []string{"1"},
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        strconv.Itoa(int(user.ID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtKeyword))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (a authService) ValidateUser(ctx context.Context, user dto.AuthUserDto) error {
	_, err := a.repo.CheckAuth(ctx, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("login and password dosen't match: %w", err)
	}
	return nil
}
