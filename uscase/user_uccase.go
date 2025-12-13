package uscase

import (
	"Cidna/model"
	"Cidna/repository"
)

type UserUsecase interface {
	SaveUser(userUs *model.User) error
}

var UserUse *UserUsecaseImpl

type UserUsecaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(UserRepository repository.UserRepository) *UserUsecaseImpl {
	return &UserUsecaseImpl{UserRepository}
}

func (u *UserUsecaseImpl) SaveUser(userUs *model.User) error {
	return u.userRepository.Save(userUs)
}
