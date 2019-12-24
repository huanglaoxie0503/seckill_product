package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"seckill_product/common"
	"seckill_product/fronted/middlerware"
	"seckill_product/fronted/web/controllers"
	"seckill_product/repositories"
	"seckill_product/services"
)

func main() {
	// 1. 创建 Iris 实例
	app := iris.New()
	// 2. 设置错误模式，在 mvc 模式下提示错误
	app.Logger().SetLevel("debug")
	// 3. 注册模板
	template := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4. 设置模板目录
	app.StaticWeb("/public", "./fronted/web/public")
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

	//sess := sessions.New(sessions.Config{
	//	Cookie:  "hello word",
	//	Expires: 60 * time.Minute,
	//})

	// 注册控制器
	user := repositories.NewUserManagerRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	product := repositories.NewProductManager("product", db)
	order := repositories.NewOrderManagerRepository("order", db)

	productService := services.NewProductService(product)
	orderService := services.NewOrderService(order)

	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	// 用户登录验证
	proProduct.Use(middlerware.AuthControllerProduct)
	pro.Register(productService, orderService)
	pro.Handle(new(controllers.ProductController))

	_ = app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
