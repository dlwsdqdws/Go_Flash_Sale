package datamodels

type Order struct {
	ID          int64 `json:"id" sql:"ID" imooc:"ID"`
	UserId      int64 `json:"UserID" sql:"userID" imooc:"UserID"`
	ProductId   int64 `json:"ProductId" sql:"productId" imooc:"ProductId"`
	OrderStatus int64 `json:"OrderStatus" sql:"orderStatus" imooc:"OrderStatus"`
}

const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)
