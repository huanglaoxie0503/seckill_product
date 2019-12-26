package config

const (
	// 设置集群地址
	HostA = "127.0.0.1"
	HostB = "127.0.0.1"
	// 数据库IP
	Host = "127.0.0.1"
	User = "root"
	Pw   = "root0503"
	Db   = "product_info"
	Port = "3306"
	// rabbitMQ连接信息
	MQUrl = "amqp://oscar:oscar@" + HostA + ":5672/product"
)
