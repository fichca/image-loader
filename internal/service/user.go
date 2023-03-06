package service

import (
	"context"
	"database/sql"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/entity"
	"github.com/jinzhu/copier"
)

type userRepository interface {
	Add(ctx context.Context, user entity.User) error
	GetById(ctx context.Context, id int) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	DeleteById(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.User, error)
}

type imageService interface {
	GetImageUrlsByUserId(ctx context.Context, userId int) ([]string, error)
}

type UserService struct {
	repo userRepository
	is   imageService
}

func NewUserService(repo userRepository, imageService imageService) *UserService {
	return &UserService{
		repo: repo,
		is:   imageService,
	}
}

func (u *UserService) Add(ctx context.Context, user dto.UserDto) error {
	return u.repo.Add(ctx, toUserEntity(user))
}

func (u *UserService) GetById(ctx context.Context, id int) (dto.UserResponse, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	urls, err := u.is.GetImageUrlsByUserId(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return toUserResponse(user, urls), err
}

func (u *UserService) Update(ctx context.Context, user dto.UserDto) error {
	return u.repo.Update(ctx, toUserEntity(user))
}

func (u *UserService) DeleteById(ctx context.Context, id int) error {
	return u.repo.DeleteById(ctx, id)
}

func (u *UserService) GetAll(ctx context.Context) ([]dto.UserDto, error) {
	users := make([]dto.UserDto, 0)

	allUsersEntity, err := u.repo.GetAll(ctx)
	if err != nil {
		return users, err
	}

	err = copier.Copy(&users, &allUsersEntity)
	if err != nil {
		return users, err
	}

	return users, nil
}

func toUserResponse(user entity.User, ImageUrls []string) dto.UserResponse {
	return dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Login:       user.Login,
		Password:    user.Password,
		Description: user.Description.String,
		ImageUrls:   ImageUrls,
	}
}

func toUserEntity(user dto.UserDto) entity.User {
	return entity.User{
		ID:       user.ID,
		Name:     user.Name,
		Login:    user.Login,
		Password: user.Password,
		Description: sql.NullString{
			String: user.Description,
			Valid:  true,
		},
	}
}
