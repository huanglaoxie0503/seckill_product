package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"seckill_product/backend/web/controllers"
	"seckill_product/common"
	"seckill_product/repositories"
	"seckill_product/services"
)

func main() {
	// 1. 创建 Iris 实例
	app := iris.New()
	// 2. 设置错误模式，在 mvc 模式下提示错误
	app.Logger().SetLevel("debug")
	// 3. 注册模板
	template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4. 设置模板目录
	app.StaticWeb("/assets", "./backend/web/assets")
	// 出现异常跳转页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		_ = ctx.View("shared/error.html")
	})
	// 连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Print(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5. 注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	// 6. 启动服务
	_ = app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
