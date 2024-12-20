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
	// 连接管理器
	connManager        wiface.IConnManager
	afterConnStartHook func(wiface.IConnection)
	beforeConnStopHook func(wiface.IConnection)
}

func NewServer(name string) wiface.IServer {
	s := &Server{
		Name:        name,
		IPVersion:   "tcp",
		IP:          global.Conf.Host,
		Port:        global.Conf.Port,
		Version:     global.Conf.Version,
		msgHandler:  NewMsgHandle(),
		connManager: NewConnManager(),
	}
	if s.Name == "" {
		s.Name = global.Conf.Name
	}
	return s
}

func (s *Server) GetConnMgr() wiface.IConnManager {
	return s.connManager
}

func (s *Server) Start() {
	SysPrintInfo(fmt.Sprintf("start server listener at ip: %s, port: %d", s.IP, s.Port))
	SysPrintInfo(fmt.Sprintf("server name: %s, host:port: %s:%d, version: %s ", s.Name, s.IP, s.Port, s.Version))

	go func() {

		//启动工作池
		s.msgHandler.StartWorkerPool()
		//监听服务地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			SysPrintError("resolve tcp addr err: ", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			SysPrintError("listen ", s.IPVersion, "err", err)
			return
		}

		//监听成功
		SysPrintInfo(fmt.Sprintf("listen tcp success, ip_version: %v, addr: %v", s.IPVersion, addr))
		cid := uint32(0)
		for {
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				SysPrintError("accept tcp error: ", err)
				continue
			}

			if s.connManager.Len() >= global.Conf.MaxConn {
				_ = tcpConn.Close()
				SysPrintError(fmt.Sprintf("accept tcp reached max: %v", global.Conf.MaxConn))
				continue
			}

			SysPrintInfo(fmt.Sprintf("connid: %v", cid))
			dealConn := NewConnection(s, tcpConn, cid, s.msgHandler)
			cid++
			SysPrintInfo(fmt.Sprintf("next connid: %v", cid))
			Print(333)
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	SysPrintInfo("server stop")
	s.GetConnMgr().ClearConn()
}

func (s *Server) AddRouter(msgId uint32, router wiface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) Serve() {
	s.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		s.Stop()
		SysPrintInfo("exit")
	}
}

func (s *Server) SetAfterConnStart(hookFunc func(wiface.IConnection)) {
	s.afterConnStartHook = hookFunc
}

func (s *Server) SetBeforeConnStop(hookFunc func(wiface.IConnection)) {
	s.beforeConnStopHook = hookFunc
}

func (s *Server) CallAfterConnStart(conn wiface.IConnection) {
	if s.afterConnStartHook != nil {
		s.afterConnStartHook(conn)
	}
}

func (s *Server) CallBeforeConnStop(conn wiface.IConnection) {
	if s.beforeConnStopHook != nil {
		s.beforeConnStopHook(conn)
	}
}
