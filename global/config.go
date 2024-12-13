package global

import (
	"winx/wiface"
)

type Config struct {
	TcpServer        wiface.IServer
	Name             string `json:"name"` //当前服务的名称
	Host             string `json:"host"`
	Port             int    `json:"port"`
	Version          string `json:"version"`             //当前框架的版本
	MaxPacketSize    uint32 `json:"max_packet_size"`     //数据包大小
	MaxConn          int    `json:"max_conn"`            //最大连接数
	WorkerPoolSize   uint32 `json:"worker_pool_size"`    //工作池的数量
	MaxWorkerTaskLen uint32 `json:"max_worker_task_len"` //worker对应的队列的最大长度
}

type Options func(config *Config)

func WithName(name string) Options {
	return func(o *Config) {
		o.Name = name
	}
}

func WithHost(host string) Options {
	return func(o *Config) {
		o.Host = host
	}
}

func WithPort(port int) Options {
	return func(o *Config) {
		o.Port = port
	}
}

func WithVersion(v string) Options {
	return func(o *Config) {
		o.Version = v
	}
}

func WithMaxPacketSize(size uint32) Options {
	return func(o *Config) {
		o.MaxPacketSize = size
	}
}

func WithMaxConn(max int) Options {
	return func(o *Config) {
		o.MaxConn = max
	}
}

func WithWorkerPoolSize(size uint32) Options {
	return func(o *Config) {
		o.WorkerPoolSize = size
	}
}

func WithMaxWorkerTaskLen(size uint32) Options {
	return func(o *Config) {
		o.MaxWorkerTaskLen = size
	}
}

func NewOption(opts ...Options) *Config {
	o := &Config{}
	opts = append(defaultOptions(), opts...)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func defaultOptions() []Options {
	return []Options{
		WithName("winx"),
		WithHost("127.0.0.1"),
		WithPort(8888),
		WithVersion("v0.1"),
		WithMaxPacketSize(4096),
		WithMaxConn(20000),
		WithWorkerPoolSize(10),
		WithMaxWorkerTaskLen(1024),
	}
}

func NewDefault() *Config {
	return NewOption()
}

var Conf *Config

func init() {
	Conf = NewDefault()
}
