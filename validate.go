package main

import (
	"errors"
	"fmt"
	"net/http"
	"seckill_product/common"
	"seckill_product/encrypt"
)

func Check(rw http.ResponseWriter, r *http.Request) {
	// 执行正常的业务逻辑
	fmt.Println("执行Check")
}

// 统一验证拦截器，每一个接口都需提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

// 身份校验函数
func CheckUserInfo(r *http.Request) error {
	// 获取uid ,cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户UID Cookie 获取失败！")
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie 获取失败！")
	}
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("用户加密串被篡改！")
	}
	// 结果比对
	fmt.Println("结果比对")
	fmt.Println("用户ID：" + uidCookie.Value)
	fmt.Println("解密后用户ID：" + string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败")
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

// 分布式身份权限验证
func main() {
	// 1.过滤器
	filter := common.NewFilter()
	// 注册拦截器
	filter.RegisterFilter("/check", Auth)
	// 2.启动服务
	http.HandleFunc("/check", filter.Handle(Check))

	_ = http.ListenAndServe(":8080", nil)

}
