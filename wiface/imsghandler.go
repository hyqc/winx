package wiface

// IMsgHandler 消息管理抽象层
type IMsgHandler interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgId uint32, router IRouter) //添加路由
	StartWorkerPool()                       //开启工作池
	SendMsgToTaskQueue(request IRequest)    //将消息交给TaskQueue,由worker进行处理
}
