package wnet

import (
	"fmt"
	"io"
	"net"
	"time"
)

func MockClient(network, addr string) {
	fmt.Println("[CLIENT] [INFO]MockClient")
	time.Sleep(time.Second * 5)
	conn, err := net.Dial(network, addr)
	if err != nil {
		fmt.Println("[CLIENT] [ERROR] dial failed, err: ", err)
		return
	}
	defer conn.Close()
	for {
		dp := NewDataPack()
		msg, err := dp.Pack(NewMessage(1, []byte("hello world")))
		if err != nil {
			fmt.Println("[CLIENT] [ERROR] write failed, err: ", err)
			return
		}
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("[CLIENT] [ERROR] write failed, err: ", err)
			return
		}

		headerData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headerData); err != nil {
			fmt.Println("[CLIENT] [ERROR] read failed, err: ", err)
			return
		}

		msgHead, err := dp.Unpack(headerData)
		if err != nil {
			fmt.Println("[CLIENT] [ERROR] unpack failed, err: ", err)
			return
		}
		msgData := msgHead.(*Message)
		if msgHead.GetDataLen() > 0 {
			msgData.Data = make([]byte, msgHead.GetDataLen())
			if _, err := io.ReadFull(conn, msgData.Data); err != nil {
				fmt.Println("[CLIENT] [ERROR] read failed, err: ", err)
				return
			}
		}

		fmt.Println(fmt.Sprintf("[CLIENT] [INFO] read, msgId: %v, msgData: %v ", msgData.GetMsgID(), string(msgData.GetData())))
		time.Sleep(time.Second)
	}
}
