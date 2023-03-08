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
	GetUserByLoginAndPassword(ctx context.Context, login, password string) (entity.User, error)
}

type tgAuthRepo interface {
	Authorize(ctx context.Context, userID int, telegramID int64) error
	CheckTgAuth(ctx context.Context, tgID int64) (int, error)
}

type AuthService struct {
	userRepo   authRepository
	tgAuthRepo tgAuthRepo
	jwtKeyword string
}

func NewAuthService(userRepo authRepository, tgAuthRepo tgAuthRepo, jwtKeyword string) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtKeyword: jwtKeyword,
		tgAuthRepo: tgAuthRepo,
	}
}

func (a *AuthService) Authorize(ctx context.Context, login, password string) (string, error) {
	user, err := a.userRepo.GetUserByLoginAndPassword(ctx, login, password)
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

func (a *AuthService) ValidateUser(ctx context.Context, authUser dto.AuthUserDto) (userID int64, err error) {
	user, err := a.userRepo.GetUserByLoginAndPassword(ctx, authUser.Login, authUser.Password)
	if err != nil {
		return 0, fmt.Errorf("login and password dosen't match: %w", err)
	}
	return user.ID, nil
}

func (a *AuthService) AuthorizeTG(ctx context.Context, tgID int64, login, password string) error {
	user, err := a.userRepo.GetUserByLoginAndPassword(ctx, login, password)
	if err != nil {
		return err
	}

	_, err = a.ValidateTGUser(ctx, tgID)
	if err == nil {
		return err
	}

	err = a.tgAuthRepo.Authorize(ctx, int(user.ID), tgID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ValidateTGUser(ctx context.Context, tgID int64) (userID int, err error) {
	userId, err := a.tgAuthRepo.CheckTgAuth(ctx, tgID)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
