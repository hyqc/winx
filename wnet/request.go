package wnet

import "winx/wiface"

type Request struct {
	conn wiface.IConnection
	data wiface.IMessage
}

func NewRequest(conn wiface.IConnection, data wiface.IMessage) wiface.IRequest {
	return &Request{
		conn: conn,
		data: data,
	}
}

func (r *Request) GetConnection() wiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.data.GetMsgID()
}
