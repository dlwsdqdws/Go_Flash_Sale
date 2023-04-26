package repositories

import (
	"errors"
	"pro-iris/datamodels"

	"gorm.io/gorm"
)

type IGormUserRepository interface {
	Select(userName string) (*datamodels.User, error)
	Insert(user *datamodels.User) (int64, error)
	SelectByID(userId int64) (*datamodels.User, error)
}

func NewGormUserManagerRepository(db *gorm.DB) IGormUserRepository {
	return &GormUserManagerRepository{db}
}

type GormUserManagerRepository struct {
	db *gorm.DB
}

func (u *GormUserManagerRepository) Select(userName string) (*datamodels.User, error) {
	user := &datamodels.User{}
	result := u.db.Where("user_name = ?", userName).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user doesn't exist")
	} else if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *GormUserManagerRepository) Insert(user *datamodels.User) (int64, error) {
	result := u.db.Create(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (u *GormUserManagerRepository) SelectByID(userId int64) (*datamodels.User, error) {
	user := &datamodels.User{}
	result := u.db.First(user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user doesn't exist")
	} else if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
