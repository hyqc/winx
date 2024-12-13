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
	fmt.Println(fmt.Sprintf("[SERVER] [INFO] PreHandle call, msgId: %v, msgData: %v", request.GetMsgID(), string(request.GetData())))
}

func (p *PingRouter) Handle(request wiface.IRequest) {
	fmt.Println(fmt.Sprintf("[SERVER] [INFO] Handle call: 2"))
	if err := request.GetConnection().SendMsg(request.GetMsgID(), request.GetData()); err != nil {
		fmt.Println("[SERVER] [ERROR] write to client failed, err: ", err)
	}
}

func (p *PingRouter) PostHandle(request wiface.IRequest) {
	fmt.Println(fmt.Sprintf("[SERVER] [INFO] PostHandle call: 3"))
}
