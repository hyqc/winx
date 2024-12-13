package wnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"winx/global"
	"winx/wiface"
)

type Connection struct {
	//当前连接
	Conn *net.TCPConn
	// 当前连接的ID号，全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	//通知退出/停止的channel
	ExitBuffChan chan bool
	//处理方法
	MsgHandler wiface.IMsgHandler
	//消息队列，分离读写消息业务
	msgChan chan []byte
}

// NewConnection 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, handler wiface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		MsgHandler:   handler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}
	return c
}

// StartReader 处理读数据
func (c *Connection) StartReader() {
	fmt.Println("[SERVER] [INFO] Reader Goroutine is running")
	defer fmt.Println("[SERVER] [INFO] ConnID: ", c.ConnID, " Reader is exit, remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()
	for {
		//拆包
		dp := NewDataPack()
		// 读取客户端消息header
		headerData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetConnection(), headerData); err != nil {
			fmt.Println("[SERVER] [ERROR] read header error:", err)
			c.ExitBuffChan <- true
			continue
		}
		fmt.Println("[SERVER] [INFO] unpack header: ", string(headerData))
		//拆解包头
		msg, err := dp.Unpack(headerData)
		if err != nil {
			fmt.Println("[SERVER] [ERROR] unpack error:", err)
			c.ExitBuffChan <- true
			continue
		}

		//拆解body
		var body []byte
		if msg.GetDataLen() > 0 {
			body = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetConnection(), body); err != nil {
				fmt.Println("[SERVER] [ERROR] read body error:", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(body)

		fmt.Println("[SERVER] [INFO] read from client msgId: ", msg.GetMsgID())
		fmt.Println("[SERVER] [INFO] read from client msgDataLen: ", msg.GetDataLen())
		fmt.Println("[SERVER] [INFO] read from client msgData: ", string(msg.GetData()))
		req := NewRequest(c, msg)
		if global.Conf.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWriter() {
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[SERVER] [ERROR] send msg error:", err)
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

// Stop 停止连接
func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	c.Conn.Close()
	c.ExitBuffChan <- true
	close(c.ExitBuffChan)
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
	select {
	case <-c.ExitBuffChan:
		fmt.Println("[SERVER] [INFO] connection exit!")
		return
	}
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

func (c *Connection) SendMsg(msgId uint32, body []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, body))
	if err != nil {
		fmt.Println("[SERVER] [ERROR] pack error:", err)
		return err
	}
	// 写回客户端
	c.msgChan <- msg
	return nil
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
