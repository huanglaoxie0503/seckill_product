package services

import (
	"seckill_product/datamodels"
	"seckill_product/repositories"
)

// 订单服务接口
type IOrderService interface {
	GetOrderByID(orderID int64) (order *datamodels.Order, err error)
	DeleteOrderByID(orderID int64) bool
	UpdateOrder(order *datamodels.Order) (err error)
	InsertOrder(order *datamodels.Order) (orderID int64, err error)
	GetAllOrder() (orderArray []*datamodels.Order, err error)
	GetAllOrderInfo() (orderMap map[int]map[string]string, err error)
}

// 订单服务结构体
type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

// 订单服务构造函数
func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository: repository}
}

// 查询订单
func (o OrderService) GetOrderByID(orderID int64) (order *datamodels.Order, err error) {
	return o.OrderRepository.SelectByKey(orderID)
}

// 删除订单
func (o OrderService) DeleteOrderByID(orderID int64) (isOK bool) {
	return o.OrderRepository.Delete(orderID)
}

// 更新订单
func (o OrderService) UpdateOrder(order *datamodels.Order) (err error) {
	return o.OrderRepository.Update(order)
}

// 插入订单
func (o OrderService) InsertOrder(order *datamodels.Order) (orderID int64, err error) {
	return o.OrderRepository.Insert(order)
}

// 查询所有订单
func (o OrderService) GetAllOrder() (orderArray []*datamodels.Order, err error) {
	return o.OrderRepository.SelectAll()
}

// 获取所有商品关联信息
func (o OrderService) GetAllOrderInfo() (orderMap map[int]map[string]string, err error) {
	return o.OrderRepository.SelectAllWithInfo()
}
