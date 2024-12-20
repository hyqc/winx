package wiface

type IConnManager interface {
	Add(conn IConnection)                   //添加链接
	Remove(conn IConnection)                //删除链接
	Get(connID uint32) (IConnection, error) //根据ConnID获取链接
	Len() int                               //获取当前连接总数
	ClearConn()                             //清除并停止所有链接
}
