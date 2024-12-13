package wiface

// IServer 定义服务器接口
type IServer interface {
	// Start 启动服务
	Start()
	// Serve 开启业务服务
	Serve()
	// Stop 停止服务
	Stop()
	// AddRouter 添加路由
	AddRouter(msgId uint32, router IRouter)
}
