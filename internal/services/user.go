package services

import (
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/entity"
)

type repository interface {
	AddUser(user entity.User) error
	GetUser(id int) (entity.User, error)
}
type userService struct {
	repo repository
}

func NewUserService(repo repository) *userService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) AddUser(user dto.UserDto) error {
	return u.repo.AddUser(toUserEntity(user))
}

func (u *userService) GetUser(id int) (dto.UserDto, error) {
	user, err := u.repo.GetUser(id)
	return toUserDto(user), err
}

func toUserDto(user entity.User) dto.UserDto {
	return dto.UserDto{
		ID:   user.ID,
		Name: user.Name,
	}
}

func toUserEntity(user dto.UserDto) entity.User {
	return entity.User{
		ID:   user.ID,
		Name: user.Name,
	}
}
