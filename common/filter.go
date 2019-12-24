package common

import "net/http"

// 声明一个新打数据类型(函数类型)
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// 声明结构体
type Filter struct {
	// 用来存储需要拦截的Url
	filterMap map[string]FilterHandle
}

// 构造函数 Filter初始化函数
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

// 注册拦截器
func (f *Filter) RegisterFilter(url string, handler FilterHandle) {
	f.filterMap[url] = handler
}

// 根据Url 获取对应打handler
func (f *Filter) GetFilterHandle(url string) FilterHandle {
	return f.filterMap[url]
}

// 声明一个新的函数类型
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// 执行拦截器，返回函数类型
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if path == r.RequestURI {
				// 执行拦截业务逻辑
				err := handle(rw, r)
				if err != nil {
					_, _ = rw.Write([]byte(err.Error()))
					return
				}
				// 跳出循环
				break
			}
		}
		// 执行正常注册的函数
		webHandle(rw, r)
	}
}
