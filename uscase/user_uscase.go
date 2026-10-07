package uscase

import (
	"NobatApi/model"
	"NobatApi/repository"
)

type UserUsecase interface {
	SaveUser(userUs *model.User) error
	GetUserByUserName(userName string) (*model.User, error)
}

var UserUse *userUsecaseImpl

type userUsecaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(UserRepository repository.UserRepository) *userUsecaseImpl {
	return &userUsecaseImpl{UserRepository}
}

func (u *userUsecaseImpl) SaveUser(userUs *model.User) error {
	// if userUs.Name == "" {
	// 	return model.ErrUserNotFound
	// }
	return u.userRepository.Save(userUs)
}

func (u *userUsecaseImpl) GetUserByUserName(userName string) (*model.User, error) {

	return u.userRepository.FindByUsername(userName)
}
