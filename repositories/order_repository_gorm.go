package repositories

import (
	"gorm.io/gorm"
	"pro-iris/common"
	"pro-iris/datamodels"
)

type IGormOrderRepository interface {
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

func NewGormOrderManager(db *gorm.DB) IGormOrderRepository {
	return &GormOrderManager{db}
}

type GormOrderManager struct {
	db *gorm.DB
}

func (o *GormOrderManager) Insert(order *datamodels.Order) (productID int64, err error) {
	err = o.db.Create(order).Error
	return order.ID, err
}

func (o *GormOrderManager) Delete(orderID int64) (check bool) {
	err := o.db.Delete(&datamodels.Order{ID: orderID}).Error
	if err != nil {
		return false
	}
	return true
}

func (o *GormOrderManager) Update(order *datamodels.Order) (err error) {
	err = o.db.Save(order).Error
	return err
}

func (o *GormOrderManager) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	order = &datamodels.Order{}
	err = o.db.First(order, orderID).Error
	if err != nil {
		return &datamodels.Order{}, err
	}
	return order, nil
}

func (o *GormOrderManager) SelectAll() (orderArray []*datamodels.Order, err error) {
	err = o.db.Find(&orderArray).Error
	return orderArray, err
}

func (o *GormOrderManager) SelectAllWithInfo() (orderMap map[int]map[string]string, err error) {
	rows, err := o.db.Table("order").
		Select("order.ID, p.product_name, order.orderStatus").
		Joins("left join products as p on order.productId = p.id").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return common.GetResultRows(rows), err
}
