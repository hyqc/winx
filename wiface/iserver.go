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
	// GetConnMgr 获取链接管理
	GetConnMgr() IConnManager
	// SetAfterConnStart 在连接启动后调用
	SetAfterConnStart(func(IConnection))
	// SetBeforeConnStop 在连接停止后调用
	SetBeforeConnStop(func(IConnection))
	// CallAfterConnStart 在连接启动之后调用
	CallAfterConnStart(conn IConnection)
	// CallBeforeConnStop 在连接停止之前调用
	CallBeforeConnStop(conn IConnection)
}
