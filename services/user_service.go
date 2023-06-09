package services

import (
	"pro-iris/datamodels"
	"pro-iris/repositories"
	"pro-iris/tool"
)

type IUserService interface {
	IsPswSuccess(userName string, psw string) (user *datamodels.User, check bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPswSuccess(userName string, psw string) (user *datamodels.User, check bool) {
	var err error
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	check, _ = tool.ValidatePassword(psw, user.HashPassword)
	if !check {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := tool.GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}
