package datamodels

// 简单消息结构体
type Message struct {
	ProductID int64
	UserID    int64
}

// 构造函数
func NewMessage(userId int64, productId int64) *Message {
	return &Message{
		UserID:    userId,
		ProductID: productId,
	}
}
