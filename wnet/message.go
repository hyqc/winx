package wnet

import "winx/wiface"

// Message 采用Type-Len-Data(TLV)解决TCP粘包问题
// 消息结构：dataLen-msgId-data
// 消息格式：header-body
// Header格式：dataLen-msgId
// Body格式：data
// 1. 先读取固定长度的dataLen，获取消息ID和数据的总长度
// 2. 再读取dataLen长度的数据，获取真正的数据
type Message struct {
	ID      uint32
	DataLen uint32
	Data    []byte
}

func NewMessage(id uint32, data []byte) wiface.IMessage {
	return &Message{
		ID:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

func (m *Message) SetDataLen(n uint32) {
	m.DataLen = n
}
