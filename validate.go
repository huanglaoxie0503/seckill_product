package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"seckill_product/common"
	"seckill_product/encrypt"
	"strconv"
	"sync"
)

// 设置集群地址，最好内网IP
var hostArray = []string{"127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1"}
var localHost = "127.0.0.1"
var port = "8081"
var hashConsistent *common.Consistent

// 用来存放控制信息
type AccessController struct {
	// 用来存放用户想要存放打信息
	sourceArray map[int]interface{}
	sync.RWMutex
}

// 创建全局变量
//var accessController = &AccessController{sourceArray:make(map[int]interface{})}

// 获取制定的数据
func (a *AccessController) GetNewRecord(uid int) interface{} {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	data := a.sourceArray[uid]
	return data
}

// 设置记录
func (a *AccessController) SetNewRecord(uid int) {
	a.RWMutex.Lock()
	defer a.RWMutex.Unlock()
	a.sourceArray[uid] = "oscar"
}

// 分布式
func (a AccessController) GetDistributedRight(req *http.Request) bool {
	// 获取用户ID
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	// 采用一致性hash算法， 很据用户ID，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	// 判断是否是本机
	if hostRequest == localHost {
		// 执行本机校验
		return a.GetDataFromMap(uid.Value)
	} else {
		// 不是本机，充当代理访问数据返回结果
		return GetDataFromOtherMap(hostRequest, req)
	}
}

// 获取本机Map , 并且处理业务逻辑，返回的结果类型为bool类型
func (a *AccessController) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := a.GetNewRecord(uidInt)
	// 执行逻辑
	if data != nil {
		return true
	}
	return false
}

// 获取其他节点处理结果
func GetDataFromOtherMap(host string, request *http.Request) bool {
	// 获取uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}
	// 获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}
	// 模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/access", nil)
	if err != nil {
		return false
	}
	// 手动指定, 排查多余cookies
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	// 添加 cookie 到模拟请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	// 获取返回结果
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	// 判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

// 执行正常的业务逻辑
func Check(rw http.ResponseWriter, r *http.Request) {
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
		return errors.New("用户Cookie 的 uid 获取失败！")
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
	// 负载均衡器设置,采用一致性哈希算法
	hashConsistent = common.NewConsistent()
	// 采用一致性哈希算法添加节点
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	// 1.过滤器
	filter := common.NewFilter()
	// 注册拦截器
	filter.RegisterFilter("/check", Auth)
	// 2.启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	// 设置监听端口
	_ = http.ListenAndServe(":8080", nil)
}
