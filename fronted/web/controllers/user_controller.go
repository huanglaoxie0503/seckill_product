package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"seckill_product/datamodels"
	"seckill_product/encrypt"
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
	// 1. 获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)
	// 2. 验证账号密码
	user, isOk := c.Service.IsPwdSuccess(userName, password)
	if !isOk {
		return mvc.Response{
			Path: "/user/login",
		}
	}
	// 3. 写入用户ID到cookie中
	tool.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidStr, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}
	// 写入浏览器Cookies
	tool.GlobalCookie(c.Ctx, "sign", uidStr)

	c.Session.Set("userID", strconv.FormatInt(user.ID, 10))
	return mvc.Response{
		Path: "/product/",
	}
}
