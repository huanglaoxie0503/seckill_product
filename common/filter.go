package common

import "net/http"

// 声明一个新打数据类型(函数类型)
type FilterHandler func(rw http.ResponseWriter, req *http.Request) error

// 声明结构体
type Filter struct {
	filterMap map[string]FilterHandler
}

// 初始化函数
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandler)}
}

// 注册拦截器
func (f *Filter) RegisterFilter(url string, handler FilterHandler) {
	f.filterMap[url] = handler
}

// 根据Url 获取对应打handler
func (f *Filter) GetFilterHandler(url string) FilterHandler {
	return f.filterMap[url]
}

// 声明一个新的函数类型
type WebHandler func(rw http.ResponseWriter, req *http.Request)

// 执行拦截器，返回函数类型
func (f *Filter) Handler(webHandler WebHandler) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handler := range f.filterMap {
			if path == r.RequestURI {
				// 执行拦截业务逻辑
				err := handler(rw, r)
				if err != nil {
					_, _ = rw.Write([]byte(err.Error()))
					return
				}
				// 跳出循环
				break
			}
		}
		// 执行正常注册的函数
		webHandler(rw, r)
	}
}
