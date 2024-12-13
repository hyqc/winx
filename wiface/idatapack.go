package wiface

// IDataPack 定义封包解包接口
type IDataPack interface {
	GetHeadLen() uint32                    //获取包头长度
	Pack(message IMessage) ([]byte, error) //封包
	Unpack(body []byte) (IMessage, error)  //解包
}
