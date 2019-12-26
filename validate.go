package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"seckill_product/common"
	"seckill_product/config"
	"seckill_product/datamodels"
	"seckill_product/encrypt"
	"seckill_product/rabbitmq"
	"strconv"
	"sync"
)

// 设置集群地址，最好内网IP,可以手动指定
var hostArray = []string{config.HostA, config.HostB}
var localHost = ""

// 数量控制接口服务器内网IP，或者GetOne的SLB内网IP
var GetOneIp = "127.0.0.1"
var GetOnePort = "8083"
var port = "8083"
var hashConsistent *common.Consistent
var rabbitMQValidate *rabbitmq.RabbitMQ

// 用来存放控制信息
type AccessController struct {
	// 用来存放用户想要存放打信息
	sourceArray map[int]interface{}
	sync.RWMutex
}

// 创建全局变量
var accessController = &AccessController{sourceArray: make(map[int]interface{})}

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
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, request)
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

// 模拟请求
func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	//获取Uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	//获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	//模拟接口访问，
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}

	//手动指定，排查多余cookies
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	//获取返回结果
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessController.GetDistributedRight(r)
	if !right {
		_, _ = w.Write([]byte("false"))
		return
	}
	_, _ = w.Write([]byte("true"))
	return
}

// 执行正常的业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("执行Check")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		_, _ = w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)

	// 获取用户cookies
	userCookie, err := r.Cookie("uid")
	if err != nil {
		_, _ = w.Write([]byte("false"))
		return
	}

	// 1.分布式权限验证
	right := accessController.GetDistributedRight(r)
	if right == false {
		_, _ = w.Write([]byte("false"))
	}
	// 2.获取数量控制权限，防止出现超卖
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		_, _ = w.Write([]byte("false"))
		return
	}
	// 判断数量控制接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			// 整合下单业务逻辑
			// 1.获取商品ID
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				_, _ = w.Write([]byte("false"))
				return
			}
			// 2.获取用户ID
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				_, _ = w.Write([]byte("false"))
				return
			}
			// 3.创建消息体
			message := datamodels.NewMessage(userID, productID)
			// 类型转换
			byteMessage, err := json.Marshal(message)
			if err != nil {
				_, _ = w.Write([]byte("false"))
				return
			}
			// 4.生产消息
			err = rabbitMQValidate.PublishSimple(string(byteMessage))
			if err != nil {
				_, _ = w.Write([]byte("false"))
				return
			}
			_, _ = w.Write([]byte("true"))
			return
		}
	}
	_, _ = w.Write([]byte("false"))
	return
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
	// 自动获取本机IP
	localIP, err := common.GetNativeIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIP
	fmt.Println(localHost)

	// 创建 rabbitMq
	rabbitMQValidate = rabbitmq.NewRabbitMQSimple("OscarProduct")
	defer rabbitMQValidate.Destroy()

	// 1.过滤器
	filter := common.NewFilter()
	// 注册拦截器
	filter.RegisterFilter("/check", Auth)
	filter.RegisterFilter("checkRight", Auth)
	// 2.启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("checkRight", filter.Handle(CheckRight))
	// 设置监听端口
	_ = http.ListenAndServe(":8080", nil)
}
