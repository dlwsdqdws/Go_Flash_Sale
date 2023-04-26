package repositories

import (
	"gorm.io/gorm"
	"pro-iris/datamodels"
)

// IGormProduct 1. define interface
type IGormProduct interface {
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
	SubProductNum(productID int64) error
}

// GormProductManager 2. implement interface
type GormProductManager struct {
	db *gorm.DB
}

func NewGormProductManager(db *gorm.DB) IGormProduct {
	return &GormProductManager{db: db}
}

func (p *GormProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	// 1. create a new product record
	result := p.db.Create(product)
	if result.Error != nil {
		return 0, result.Error
	}

	// 2. return the newly created product's ID
	return product.ID, nil
}

func (p *GormProductManager) Delete(productId int64) bool {
	// 1. delete the product record with the given ID
	result := p.db.Delete(&datamodels.Product{}, productId)

	// 2. check if any records were affected
	return result.RowsAffected > 0
}

func (p *GormProductManager) Update(product *datamodels.Product) error {
	// 1. update the product record
	result := p.db.Save(product)

	// 2. check for any errors
	if result.Error != nil {
		return result.Error
	}

	// 3. check if any records were affected
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *GormProductManager) SelectByKey(productId int64) (productRes *datamodels.Product, err error) {
	// 1. retrieve the product record with the given ID
	result := p.db.First(&productRes, productId)

	// 2. check for any errors
	if result.Error != nil {
		return nil, result.Error
	}

	// 3. check if any records were found
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return productRes, nil
}

func (p *GormProductManager) SelectAll() (productArray []*datamodels.Product, errRes error) {
	// 1. retrieve all product records

	result := p.db.Find(&productArray)
	// 2. check for any errors
	if result.Error != nil {
		return nil, result.Error
	}

	return productArray, nil
}

func (p *GormProductManager) SubProductNum(productID int64) error {
	// 1. decrement the productNum column for the record with the given ID
	result := p.db.Model(&datamodels.Product{}).Where("id = ?", productID).Update("product_num", gorm.Expr("product_num - ?", 1))

	// 2. check for any errors
	if result.Error != nil {
		return result.Error
	}

	// 3. check if any records were affected
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
