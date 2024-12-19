package wiface

import "net"

// IConnection 定义连接接口
type IConnection interface {
	// Start 开始连接
	Start()
	// Stop 停止当前连接
	Stop()
	// GetConnID 获取当前连接ID
	GetConnID() uint32
	// GetConnection 获取连接
	GetConnection() net.Conn
	// SendMsg 发包
	SendMsg(msgId uint32, body []byte) error
	// SendBuffMsg 发送带缓冲的消息
	SendBuffMsg(msgId uint32, data []byte) error
}
