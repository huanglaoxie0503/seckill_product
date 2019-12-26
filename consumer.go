package main

import (
	"fmt"
	"seckill_product/common"
	"seckill_product/rabbitmq"
	"seckill_product/repositories"
	"seckill_product/services"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	//创建product数据库操作实例
	product := repositories.NewProductManager("product", db)
	// 创建product service
	productService := services.NewProductService(product)
	// 创建 Order 数据库实例
	order := repositories.NewOrderManagerRepository("order", db)
	// 创建 Order Service
	orderService := services.NewOrderService(order)
	// RabbitMQ 消费消息
	rabbitMQConsumerSimple := rabbitmq.NewRabbitMQSimple("OscarProduct")
	rabbitMQConsumerSimple.ConsumeSimple(orderService, productService)
}
