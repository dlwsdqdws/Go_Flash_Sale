package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"pro-iris/datamodels"
	"pro-iris/repositories"
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
	check, _ = ValidatePassword(psw, user.HashPassword)
	if !check {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPsw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPsw), bcrypt.DefaultCost)
}

func ValidatePassword(userPsw string, hashedPsw string) (check bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPsw), []byte(userPsw)); err != nil {
		return false, errors.New("Incorrect Password!")
	}
	return true, nil
}
