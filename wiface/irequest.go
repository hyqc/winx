package wiface

type IRequest interface {
	GetConnection() IConnection //获取请求的连接
	GetData() []byte            //获取请求消息的数据
	GetMsgID() uint32
}
