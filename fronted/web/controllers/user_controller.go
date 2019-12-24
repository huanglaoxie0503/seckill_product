package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"seckill_product/datamodels"
	"seckill_product/services"
	"seckill_product/tool"
	"strconv"
)

//
type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	Session sessions.Session
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister() {
	// 从前端获取表单数据
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)
	// ozzo-validation
	user := &datamodels.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: password,
	}
	_, err := c.Service.AddUser(user)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return
	}
	c.Ctx.Redirect("/user/login")
	return
}

// 登录

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response {
	// 从前端获取表单数据
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)
	user, isOk := c.Service.IsPwdSuccess(userName, password)
	if !isOk {
		return mvc.Response{
			Path: "/user/login",
		}
	}
	tool.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	c.Session.Set("userID", strconv.FormatInt(user.ID, 10))
	return mvc.Response{
		Path: "/product/",
	}
}
