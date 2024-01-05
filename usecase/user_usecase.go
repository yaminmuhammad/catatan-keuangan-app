package usecase

import (
	"fmt"
	"time"

	"catatan-keuangan-app/entity"
	"catatan-keuangan-app/repository"
)

type UserUseCase interface {
	RegisterNewUser(payload entity.User) (entity.User, error)
	FindUserByID(id string) (entity.User, error)
	FindUserByUsernamePassword(username, password string) (entity.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func (u *userUseCase) RegisterNewUser(payload entity.User) (entity.User, error) {
	userExist, _ := u.repo.GetByUsername(payload.Username)
	if userExist.Username == payload.Username {
		return entity.User{}, fmt.Errorf("user with username: %s already exists", payload.Username)
	}
	payload.Role = "user"
	payload.UpdatedAt = time.Now()
	return u.repo.Create(payload)
}

func (u *userUseCase) FindUserByID(id string) (entity.User, error) {
	return u.repo.Get(id)
}

func (u *userUseCase) FindUserByUsernamePassword(username, password string) (entity.User, error) {
	userExist, err := u.repo.GetByUsername(username)
	if err != nil {
		return entity.User{}, fmt.Errorf("user doesn't exists")
	}

	if userExist.Password != password {
		return entity.User{}, fmt.Errorf("user doesn't exists")
	}

	return userExist, nil
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
