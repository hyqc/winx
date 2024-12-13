package wnet

import (
	"fmt"
	"winx/global"
	"winx/wiface"
)

type MsgHandle struct {
	ApisM          map[uint32]wiface.IRouter
	WorkerPoolSize uint32                 //业务工作池worker数量
	TaskQueue      []chan wiface.IRequest //负责取消息任务的队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		ApisM:          make(map[uint32]wiface.IRouter),
		WorkerPoolSize: global.Conf.WorkerPoolSize,
		TaskQueue:      make([]chan wiface.IRequest, global.Conf.WorkerPoolSize),
	}
}

func (m *MsgHandle) DoMsgHandler(request wiface.IRequest) {
	handler, ok := m.ApisM[request.GetMsgID()]
	if !ok {
		fmt.Println("[SERVER] [ERROR] no handler for msgID: ", request.GetMsgID())
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgId uint32, router wiface.IRouter) {
	if _, ok := m.ApisM[msgId]; ok {
		panic("[SERVER] [ERROR] msgID: " + fmt.Sprintf("%d", msgId) + " has been registered")
	}
	m.ApisM[msgId] = router
	fmt.Println("[SERVER] [INFO] add router for msgID: ", msgId)
}

func (m *MsgHandle) StartWorker(workerId int, taskQueue chan wiface.IRequest) {
	fmt.Println(fmt.Sprintf("[SERVER] [INFO] worker id: %v is started", workerId))
	for {
		select {
		case req := <-taskQueue:
			m.DoMsgHandler(req)
		}
	}
}

func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan wiface.IRequest, global.Conf.MaxWorkerTaskLen)
		go m.StartWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request wiface.IRequest) {
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	m.TaskQueue[workerId] <- request
}
