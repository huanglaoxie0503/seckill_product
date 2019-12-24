package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"seckill_product/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
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
