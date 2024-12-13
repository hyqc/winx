package wnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"winx/global"
	"winx/wiface"
)

//实现封包解包

type DataPack struct {
}

var (
	ErrMsgTooLarger = errors.New("message unpack: too large msg data received")
)

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// id uint32(4字节) +  dataLen uint32(4字节)
	return 8
}

func (d *DataPack) Pack(msg wiface.IMessage) ([]byte, error) {
	//创建一个缓冲区
	dataBuff := bytes.NewBuffer([]byte{})
	//写入数据长度
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//写入消息ID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	//写入消息数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (d *DataPack) Unpack(data []byte) (wiface.IMessage, error) {
	dataBuff := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	if global.Conf.MaxPacketSize > 0 && msg.DataLen > global.Conf.MaxPacketSize {
		return nil, ErrMsgTooLarger
	}
	return msg, nil
}
