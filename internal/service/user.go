package service

import (
	"context"
	"database/sql"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/entity"
)

type repository interface {
	Add(ctx context.Context, user entity.User) error
	GetById(ctx context.Context, id int) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) error
}
type userService struct {
	repo repository
}

func (u *userService) Add(ctx context.Context, user dto.UserDto) error {
	return u.repo.Add(ctx, toUserEntity(user))
}

func (u *userService) GetById(ctx context.Context, id int) (dto.UserDto, error) {
	user, err := u.repo.GetById(ctx, id)
	return toUserDto(user), err
}

func (u *userService) Update(ctx context.Context, user dto.UserDto) error {
	//TODO implement me
	panic("implement me")
}

func (u *userService) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (u *userService) GetAll(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewUserService(repo repository) *userService {
	return &userService{
		repo: repo,
	}
}

func toUserDto(user entity.User) dto.UserDto {
	return dto.UserDto{
		ID:          user.ID,
		Name:        user.Name,
		Login:       user.Login,
		Password:    user.Password,
		Description: user.Description.String,
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
