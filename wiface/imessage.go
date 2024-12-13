package wiface

type IMessage interface {
    GetDataLen() uint32 //消息长度
    GetData() []byte //消息内容
    GetMsgID() uint32 //消息ID

	SetDataLen(uint32) //设置消息长度
	SetData([]byte) //设置消息内容
	SetMsgID(uint32) //设置消息ID
}
