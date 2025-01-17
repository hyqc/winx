package main

import (
	"fmt"
	"winx/wiface"
	"winx/wnet"
)

func main() {
	s := wnet.NewServer("winx")
	s.AddRouter(1, NewPingRouter())

	go wnet.MockClient("tcp", "127.0.0.1:8888")
	go wnet.MockClient("tcp", "127.0.0.1:8888")
	go wnet.MockClient("tcp", "127.0.0.1:8888")
	go wnet.MockClient("tcp", "127.0.0.1:8888")

	s.Serve()
}

type PingRouter struct {
	wnet.BaseRouter
	MsgId uint32
}

func NewPingRouter() *PingRouter {
	return &PingRouter{
		MsgId: 1,
	}
}

func (p *PingRouter) PreHandle(request wiface.IRequest) {
	wnet.SysPrintInfo(fmt.Sprintf("PreHandle call, msgId: %v, msgData: %v", request.GetMsgID(), string(request.GetData())))
}

func (p *PingRouter) Handle(request wiface.IRequest) {
	wnet.SysPrintInfo(fmt.Sprintf("Handle call: 2"))
	if err := request.GetConnection().SendMsg(request.GetMsgID(), request.GetData()); err != nil {
		wnet.SysPrintError("write to client failed, err: ", err)
	}
}

func (p *PingRouter) PostHandle(request wiface.IRequest) {
	wnet.SysPrintInfo(fmt.Sprintf(" PostHandle call: 3"))
}
