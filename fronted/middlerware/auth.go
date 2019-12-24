package middlerware

import "github.com/kataras/iris"

/*
	Go Iris 中间件开发
*/
func AuthControllerProduct(ctx iris.Context) {
	uid := ctx.GetCookie("uid")
	if uid == "" {
		ctx.Application().Logger().Debug("必须先登录")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("已经登录")
	// 继续执行下一个handler
	ctx.Next()
}
