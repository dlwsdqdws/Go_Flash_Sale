package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" imooc:"ID" gorm:"primaryKey;autoIncrement;not null"`
	ProductName  string `json:"ProductName" sql:"productName" imooc:"ProductName" gorm:"size:255;not null"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" imooc:"ProductNum" gorm:"not null"`
	ProductImage string `json:"ProductImage" sql:"productImage" imooc:"ProductImage" gorm:"size:255;null"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" imooc:"ProductUrl" gorm:"size:255;null"`
}
