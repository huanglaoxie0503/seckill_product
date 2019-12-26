package controllers

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"seckill_product/datamodels"
	"seckill_product/rabbitmq"
	"seckill_product/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	RabbitMQ       rabbitmq.RabbitMQ
	Session        *sessions.Session
}

func (p *ProductController) GetDetail() mvc.View {
	idStr := p.Ctx.URLParam("productID")
	// str --> int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	product, err := p.ProductService.GetProductByID(id)
	return mvc.View{
		Layout: "/shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() []byte {
	productStr := p.Ctx.URLParam("productID")
	userStr := p.Ctx.GetCookie("uid")
	productID, err := strconv.ParseInt(productStr, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	userID, err := strconv.ParseInt(userStr, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	// 创建消息体
	message := datamodels.NewMessage(userID, productID)
	// 类型转换
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	// 调用 RabbitMQ 发布消息
	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return []byte("true")
}

// 该方法适用于请求不大的场景，容易摧毁数据库
func (p *ProductController) GetOrderSimple() mvc.View {
	productStr := p.Ctx.URLParam("productID")
	userStr := p.Ctx.GetCookie("uid")
	productID, err := strconv.Atoi(productStr)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderID int64
	showMessage := "抢购失败！"
	// 判断商品数量是否满足需求
	if product.ProductNum > 0 {
		// 扣除商品数量
		product.ProductNum -= 1
		// 更新数据库
		errUpdate := p.ProductService.UpdateProduct(product)
		if errUpdate != nil {
			p.Ctx.Application().Logger().Debug(errUpdate)
		}
		// 抢购商品 创建订单
		userID, err := strconv.Atoi(userStr)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		order := &datamodels.Order{
			UserId:      int64(userID),
			ProductId:   int64(productID),
			OrderStatus: datamodels.OrderSuccess,
		}
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功！"
		}
	}
	return mvc.View{
		Layout: "/shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}
}
