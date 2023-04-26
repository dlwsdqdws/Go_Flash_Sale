package services

import (
	"pro-iris/datamodels"
	"pro-iris/repositories"
)

type IGormOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(*datamodels.Message) (int64, error)
}

func NewGormOrderService(repository repositories.IGormOrderRepository) IGormOrderService {
	return &GormOrderService{repository}
}

type GormOrderService struct {
	GormOrderRepository repositories.IGormOrderRepository
}

func (o *GormOrderService) GetOrderByID(orderID int64) (order *datamodels.Order, err error) {
	return o.GetOrderByID(orderID)
}

func (o *GormOrderService) DeleteOrderByID(orderID int64) bool {
	return o.DeleteOrderByID(orderID)
}

func (o *GormOrderService) UpdateOrder(order *datamodels.Order) error {
	return o.UpdateOrder(order)
}

func (o *GormOrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return o.InsertOrder(order)
}

func (o *GormOrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return o.GetAllOrder()
}

func (o *GormOrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.GetAllOrderInfo()
}

func (o *GormOrderService) InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error) {
	order := &datamodels.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}
