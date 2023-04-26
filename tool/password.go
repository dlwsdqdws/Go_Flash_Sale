package tool

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(userPsw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPsw), bcrypt.DefaultCost)
}

func ValidatePassword(userPsw string, hashedPsw string) (check bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPsw), []byte(userPsw)); err != nil {
		return false, errors.New("Incorrect Password!")
	}
	return true, nil
}
