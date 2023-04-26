package services

import (
	"pro-iris/datamodels"
	"pro-iris/repositories"
)

type IGormProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	SubNumberOne(productID int64) error
}

type GormProductService struct {
	productRepository repositories.IGormProduct
}

func NewGormProductService(repository repositories.IGormProduct) IGormProductService {
	return &GormProductService{repository}
}

func (p *GormProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(productID)
}

func (p *GormProductService) GetAllProduct() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

func (p *GormProductService) DeleteProductByID(productID int64) bool {
	return p.productRepository.Delete(productID)
}

func (p *GormProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

func (p *GormProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}

func (p *GormProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductNum(productID)
}
