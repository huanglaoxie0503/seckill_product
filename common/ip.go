package common

import (
	"errors"
	"net"
)

// 获取本机IP
func GetNativeIp() (string, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, adder := range address {
		// 检查IP地址判断是否回环地址
		if ipNet, isOk := adder.(*net.IPNet); isOk && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("获取IP地址异常！")
}
