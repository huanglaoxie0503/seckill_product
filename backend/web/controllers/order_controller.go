package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"seckill_product/services"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.OrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("查询订单失败！")
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}
