package wnet

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"winx/global"
	"winx/wiface"
)

type Server struct {
	// 服务名称
	Name string
	// 服务绑定的IP地址
	IP string
	// 服务绑定的端口号
	Port int
	//版本
	Version string
	// tcp4 or other
	IPVersion  string
	msgHandler wiface.IMsgHandler
}

func NewServer(name string) wiface.IServer {
	s := &Server{
		Name:       name,
		IPVersion:  "tcp",
		IP:         global.Conf.Host,
		Port:       global.Conf.Port,
		Version:    global.Conf.Version,
		msgHandler: NewMsgHandle(),
	}
	if s.Name == "" {
		s.Name = global.Conf.Name
	}
	return s
}

func (s *Server) Start() {
	fmt.Println(fmt.Sprintf("[SERVER] [INFO] start server listener at ip: %s, port: %d", s.IP, s.Port))
	fmt.Println(fmt.Sprintf("[SERVER] [INFO] server name: %s, host:port: %s:%d, version: %s ", s.Name, s.IP, s.Port, s.Version))

	go func() {
		//启动工作池
		s.msgHandler.StartWorkerPool()
		//监听服务地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("[SERVER] [ERROR] resolve tcp addr err: ", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("[SERVER] [ERROR] listen", s.IPVersion, "err", err)
			return
		}
		//监听成功
		fmt.Println(fmt.Sprintf("[SERVER] [INFO] listen tcp success, ip_version: %v, addr: %v", s.IPVersion, addr))
		cid := uint32(0)
		for {
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[SERVER] [ERROR] accept tcp error: ", err)
				continue
			}

			dealConn := NewConnection(tcpConn, cid, s.msgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) AddRouter(msgId uint32, router wiface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) Stop() {
	fmt.Println("[SERVER] [INFO] server stop")
}

func (s *Server) Serve() {
	s.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		fmt.Println("[SERVER] [INFO] exit")
	}
}
