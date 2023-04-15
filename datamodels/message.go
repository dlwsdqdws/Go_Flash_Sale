package datamodels

type Message struct {
	ProductID int64
	UserID    int64
}

func NewMessage(userID int64, productID int64) *Message {
	return &Message{
		UserID:    userID,
		ProductID: productID,
	}
}
