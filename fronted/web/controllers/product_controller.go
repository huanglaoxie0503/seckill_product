package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"seckill_product/datamodels"
	"seckill_product/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
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

func (p *ProductController) GetOrder() mvc.View {
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
